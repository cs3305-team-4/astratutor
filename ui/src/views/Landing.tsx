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
    --gradient-color-1: #00ef5c;
    --gradient-color-2: #1fad0d;
    --gradient-color-3: #6ff16a;
    --gradient-color-4: #1d7711;
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
        <Row style={{ height: 'calc(100vh - 72px)' }}>
          <Col style={{ overflow: 'hidden' }}>
            <Typography
              style={{
                position: 'absolute',
                width: '100%',
                height: '100vh',
                zIndex: 1,
                background:
                  'linear-gradient(135deg, rgba(255, 255, 255, 0) 0%, rgba(255, 255, 255, 0) 50%, rgb(27 27 27) 50%)',
              }}
            >
              <img
                src="/logo.svg"
                alt="AstraTutor"
                style={{ float: 'right', right: -120, marginTop: 60, height: 800, position: 'absolute' }}
              />
              <Title style={{ zIndex: 1000, color: '#ffffff', position: 'absolute', right: 70, top: '50%' }} level={1}>
                Take your marks to the stars!
              </Title>
            </Typography>
            <canvas
              data-js-darken-top
              data-transition-in
              id="gradient"
              style={{ position: 'relative', width: '100%', height: '100vh' }}
              ref={(r) => {
                canvasRef.current = r ?? undefined;
              }}
            />
          </Col>
        </Row>
      </Hero>
    </Layout>
  );
}
