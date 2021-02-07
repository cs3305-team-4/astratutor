import React from 'react';
import styled from 'styled-components';

import { Typography, Layout, Row, Col } from 'antd';

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const Hero = styled(Content)`
  background-color: rgba(233, 233, 233);
`;

export default function Landing() {
  return (
    <Layout>
      <Hero>
        <Row style={{ height: 'calc(100vh - 72px)' }} align="middle" justify="center">
          <Col>
            <Typography>
              <Title level={1}>Take your marks to the stars!</Title>
            </Typography>
          </Col>
        </Row>
      </Hero>
    </Layout>
  );
}
