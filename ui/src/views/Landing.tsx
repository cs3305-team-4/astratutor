import React, { useEffect, useRef } from 'react';
import styled from 'styled-components';

import { Typography, Layout, Row, Col } from 'antd';
import Gradient from '../api/gradient';

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const Hero = styled(Content)`
  background-color: rgba(233, 233, 233);
  #gradient {
    width: 100%;
    height: 100%;
    --gradient-color-1: #ef008f;
    --gradient-color-2: #6ec3f4;
    --gradient-color-3: #7038ff;
    --gradient-color-4: #ffba27;
  }
`;

export function Landing(): React.ReactElement {
  const canvasRef = useRef<HTMLCanvasElement>();

  useEffect(() => {
    const gradient = new Gradient();
    gradient.initGradient('#gradient');
  }, [canvasRef]);

  return (
    <Layout>
      <Hero>
        <Row style={{ height: 'calc(100vh - 72px)' }} align="middle" justify="center">
          <Col>
            <canvas
              data-js-darken-top
              data-transition-in
              id="gradient"
              ref={(r) => {
                canvasRef.current = r ?? undefined;
              }}
            />
            <Typography>
              <Title level={1}>Take your marks to the stars!</Title>
            </Typography>
          </Col>
        </Row>
      </Hero>
    </Layout>
  );
}
