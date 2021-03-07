import React, { useEffect, useRef } from 'react';
import styled from 'styled-components';

import { Typography, Layout, Row, Col } from 'antd';

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const Hero = styled(Content)`
  background-color: rgba(233, 233, 233);
`;

export function Landing(): React.ReactElement {
  const canvasRef = useRef<HTMLCanvasElement>();

  useEffect(() => {
    const context = canvasRef.current?.getContext('2d');
    if (!context) return;
    let time = 0;
    context.scale(5, 5);

    const color = function (x: number, y: number, r: number, g: number, b: number) {
      context.fillStyle = `rgb(${r}, ${g}, ${b})`;
      context.fillRect(x, y, 10, 10);
    };
    const R = function (x: number, y: number, time: number) {
      return Math.floor(192 + 64 * Math.cos((x * x - y * y) / 300 + time));
    };

    const G = function (x: number, y: number, time: number) {
      return Math.floor(192 + 64 * Math.sin((x * x * Math.cos(time / 4) + y * y * Math.sin(time / 3)) / 300));
    };

    const B = function (x: number, y: number, time: number) {
      return Math.floor(
        192 + 64 * Math.sin(5 * Math.sin(time / 9) + ((x - 100) * (x - 100) + (y - 100) * (y - 100)) / 1100),
      );
    };

    const startAnimation = function () {
      for (let x = 0; x <= 30; x++) {
        for (let y = 0; y <= 30; y++) {
          color(x, y, R(x, y, time), G(x, y, time), B(x, y, time));
        }
      }
      time = time + 0.06;
      window.requestAnimationFrame(startAnimation);
    };

    startAnimation();
  }, [canvasRef]);

  return (
    <Layout>
      <Hero>
        <Row style={{ height: 'calc(100vh - 72px)' }} align="middle" justify="center">
          <Col>
            <canvas
              style={{ width: '100vw', height: '100vh' }}
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
