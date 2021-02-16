import React, { ReactElement, useContext, useState } from 'react';

import { SubjectDTO, TutorSubjectsDTO } from '../api/definitions';

import { Link, useLocation, useHistory } from 'react-router-dom';

import { Typography, Layout, Card, Row, Col, List, Button, Input, Select, Space, Tabs, Tag } from 'antd';
import { useAsync } from 'react-async-hook';
import { APIContext } from '../api/api';

const { Title, Paragraph } = Typography;
const { Content } = Layout;

export function Tutors(): ReactElement {
  const api = useContext(APIContext);
  const [tutors, setTutors] = useState<TutorSubjectsDTO[] | undefined>(undefined);
  const [subjects, setSubjects] = useState<SubjectDTO[] | undefined>(undefined);
  const [filters, setFilters] = useState<string[]>([]);

  const query = new URLSearchParams(useLocation().search);
  const history = useHistory();

  useAsync(async () => {
    if (query.has('filter')) {
      const filterValues = query.get('filter').split(',');
      setFilters(filterValues);
      setTutors(await api.services.readTutors(filterValues));
    } else {
      setTutors(await api.services.readTutors());
    }
    setSubjects(await api.services.readSubjects());
  }, []);

  const onFiltersChange = async (e: string[]) => {
    history.push(e.length > 0 ? `/subjects/tutors?filter=${e.join(',')}` : '/subjects/tutors');
    setFilters(e);
    setTutors(await api.services.readTutors(e));
  };

  const onSearch = (searchVal: string) => {
    console.log(searchVal);
    // TODO: Add search functionality to /subjects/tutors endpoint
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
              <Input.Search key="2" placeholder="Search for a tutor" onSearch={onSearch} />
            </Space>
          </Row>
          <List
            itemLayout="vertical"
            size="large"
            loading={tutors === undefined}
            dataSource={tutors}
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
