import React from 'react';
import styled from 'styled-components';


import {
  Typography,
  Layout
} from "antd";

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const Hero = styled(Content)`
  padding: 25vh 0;
  text-align: center;

  h1 {
    font-size: 4rem;
  }

  background-color: rgba(233,233,233);
`

export default function Landing() {
  return (
    <Layout>
      <Hero>
        <Typography>
          <h1>
            Give your grades a blast off!
          </h1>
        </Typography>
      </Hero>
    </Layout>
  )
}