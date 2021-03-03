import React, { ReactElement } from 'react';
import styled from 'styled-components';

import { Link } from 'react-router-dom';

import { Typography, Layout, Card, Row, Col, Image } from 'antd';
import { SubjectDTO } from '../api/definitions';

const { Title, Paragraph, Text } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const StyledRow = styled(Row)<{ color: string }>`
  height: 100px;
  font-size: 1.5em;
  overflow: hidden;
  color: ${(props) => props.color};
  background-color: white;
  border-bottom: 5px solid ${(props) => props.color};
  font-weight: bolder;
  margin: 10px;
  transition: all 0.2s;
  &:hover {
    border-bottom: 10px solid ${(props) => props.color};
    font-size: 2em;
  }
`;

export interface SubjectProps {
  subject: SubjectDTO;
}

export default function Subject(props: SubjectProps): ReactElement {
  return (
    <Link to={`/subjects/tutors?filter=${props.subject.slug}`}>
      <StyledRow
        color={'#' + ('00000' + ((Math.random() * (1 << 24) - 2000) | 0).toString(16)).slice(-6)}
        align="middle"
      >
        <Col span={24} style={{ padding: '20px' }}>
          {props.subject.name}
        </Col>
        {/* <Col span={8}>
          <img height="150px" src={props.subject.image} alt={props.subject.name} />
        </Col> */}
      </StyledRow>
    </Link>
  );
}
