import React from 'react';
import styled from 'styled-components';


import {
  Typography,
  Layout
} from "antd";

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const Hero = styled(Content)`


`
export default function AccountProfile() {
  return (
    <Layout>
      <Typography>
        <Title>
          Profile
        </Title>
      </Typography>
    </Layout>
  )
}