import React, { ReactElement } from 'react';
import styled from 'styled-components';

import { Link } from 'react-router-dom';

import { Typography, Layout, Card, Row, Col, Image } from 'antd';
import { SubjectDTO } from '../api/definitions';

const { Title, Paragraph, Text } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const colors = [
  '#ef25a4',
  '#1c73cd',
  '#05c760',
  '#17eaaa',
  '#bd00ff',
  '#ff9a00',
  '#8900ff',
  '#d9534f',
  '#5cb85c',
  '#ff4d00',
  '#ffc100',
  '#00e6f9',
  '#e51894',
  '#29e518',
  '#e56b18',
];
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
    <Link to={`/subjects/tutors?filter=${props.subject.slug}&sort=featured`}>
      <StyledRow
        color={(() => {
          let s = 0;
          for (const i in props.subject.name.split('')) {
            s += props.subject.name.charCodeAt(+i);
          }
          return colors[s % (colors.length - 1)];
        })()}
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
