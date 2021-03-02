import React, { ReactElement, useContext, useState } from 'react';

import { SubjectDTO, TutorSubjectsDTO } from '../api/definitions';

import { Link, useLocation, useHistory } from 'react-router-dom';

import { Typography, Layout, Card, Row, Col, List, Button, Input, Select, Space, Tabs, Tag, Pagination } from 'antd';
import { useAsync } from 'react-async-hook';
import { APIContext } from '../api/api';

const { Title, Paragraph } = Typography;
const { Content } = Layout;

export function Tutors(): ReactElement {
  const api = useContext(APIContext);
  const [tutors, setTutors] = useState<TutorSubjectsDTO[] | undefined>(undefined);
  const [subjects, setSubjects] = useState<SubjectDTO[] | undefined>(undefined);
  const [filters, setFilters] = useState<string[]>([]);
  const [search, setSearch] = useState<string>('');
  const [searchBox, setSearchBox] = useState<string>('');

  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [totalPages, setTotalPages] = useState<number>(1);

  const query = new URLSearchParams(useLocation().search);
  const history = useHistory();

  const updatePath = (newPage: number, newPageSize: number, newFilters: string[], newQuery: string) => {
    const path = '/subjects/tutors';
    const queries: string[] = [];

    if (newFilters.length > 0) {
      queries.push(`filter=${newFilters.join(',')}`);
    }
    if (newPage > 1) {
      queries.push(`page=${newPage}`);
    }
    if (pageSize !== 10) {
      queries.push(`page_size=${pageSize}`);
    }
    if (newQuery.length > 0) {
      queries.push(`query=${newQuery}`);
    }
    if (queries.length > 0) history.push(path + '?' + queries.join('&'));
    else history.push(path);
  };

  // Initial Page Load
  useAsync(async () => {
    if (query.has('filter')) {
      setFilters((query.get('filter') ?? '').split(','));
    }
    if (query.has('page')) {
      setCurrentPage(+(query.get('page') ?? 0));
    }
    if (query.has('page_size')) {
      setPageSize(+(query.get('page_size') ?? 0));
    }
    if (query.has('query')) {
      setSearch(query.get('query') ?? '');
      setSearchBox(query.get('query') ?? '');
    }

    setSubjects(await api.services.readSubjects());
  }, []);

  // Called every tune dependencies change
  useAsync(async () => {
    const res = await api.services.readTutors(currentPage, pageSize, filters, search);
    console.log(res);

    setTotalPages(res.total_pages);
    setTutors(res.items);

    updatePath(currentPage, pageSize, filters, search);
  }, [currentPage, pageSize, filters, search]);

  const onFiltersChange = async (e: string[]) => {
    setCurrentPage(1);
    setFilters(e);
  };

  const onSearch = (searchVal: string) => {
    setCurrentPage(1);
    setSearch(searchVal);
    // TODO: Add search functionality to /subjects/tutors endpoint
  };

  const onPaginationUpdate = (newPage: number, newPageSize: number) => {
    setCurrentPage(newPage);
    setPageSize(newPageSize);
  };

  return (
    <Content style={{ padding: '2em 0' }}>
      <Row>
        <Col xl={{ offset: 4, span: 16 }} lg={{ offset: 2, span: 20 }} span={24}>
          <Row justify="space-between">
            <Title>Tutors</Title>
            <Space>
              <Select
                key="1"
                mode="multiple"
                allowClear
                value={filters}
                placeholder="Filter"
                onChange={onFiltersChange}
                style={{ minWidth: '200px' }}
              >
                {subjects?.map((subject, index) => (
                  <Select.Option key={index} value={subject.slug}>
                    {subject.name}
                  </Select.Option>
                ))}
              </Select>
              <Input.Search
                value={searchBox}
                key="2"
                placeholder="Search for a tutor"
                onChange={(e) => setSearchBox(e.currentTarget.value)}
                onSearch={onSearch}
              />
            </Space>
          </Row>
          <List
            itemLayout="vertical"
            size="large"
            loading={tutors === undefined}
            dataSource={tutors}
            pagination={{
              current: currentPage,
              pageSize: pageSize,
              pageSizeOptions: ['1', '10', '15', '25', '50', '100'],
              onChange: onPaginationUpdate,
              onShowSizeChange: onPaginationUpdate,
              total: totalPages * pageSize,
              showSizeChanger: true,
            }}
            renderItem={(tutor: TutorSubjectsDTO) => (
              <Card>
                <List.Item
                  key={tutor.id}
                  extra={<img width={200} src={tutor.avatar} alt="" />}
                  actions={[
                    <Link key="1" to={`/tutors/${tutor.id}/profile`}>
                      <Button type="primary">Visit Profile</Button>
                    </Link>,
                  ]}
                >
                  <List.Item.Meta
                    title={
                      <Link to={`/tutors/${tutor.id}/profile`}>
                        {tutor.first_name} {tutor.last_name}
                      </Link>
                    }
                    description={
                      <Tabs>
                        <Tabs.TabPane tab="Tutor Descrption">{tutor.description}</Tabs.TabPane>
                        {tutor.subjects.map((subject) => (
                          <Tabs.TabPane
                            key={subject.id}
                            tab={
                              <Tag color={filters.includes(subject.slug) ? 'blue' : ''}>
                                {subject.name} - €{subject.price}/Hour
                              </Tag>
                            }
                          >
                            {subject.description}
                          </Tabs.TabPane>
                        ))}
                      </Tabs>
                    }
                  />
                </List.Item>
              </Card>
            )}
          />
        </Col>
      </Row>
    </Content>
  );
}
