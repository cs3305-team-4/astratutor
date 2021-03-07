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
        <Row style={{ height: 'calc(25vh)', overflow: 'hidden' }}>
          <Col style={{ overflow: 'hidden' }}>
            <Typography
              style={{
                position: 'absolute',
                width: '100%',
                height: '80%',
                top: -300,
                zIndex: 1,
                background:
                  'linear-gradient(-187deg, rgba(255, 255, 255, 0) 0%, rgba(255, 255, 255, 0) 50%, rgb(240, 242, 245) 50%)',
              }}
            >
              <Title
                style={{
                  zIndex: 1000,
                  position: 'absolute',
                  top: '48%',
                  fontSize: '3em',
                  width: '100%',
                  textAlign: 'center',
                }}
                level={1}
              >
                A Complete Modern Tutoring online platform
              </Title>
            </Typography>
            <canvas
              data-js-darken-top
              data-transition-in
              id="gradient"
              style={{ position: 'relative', width: '100%', height: '80%' }}
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
