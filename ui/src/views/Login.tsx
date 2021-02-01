
import React from 'react';
import styled from 'styled-components';


import {
  Layout,
  PageHeader
} from "antd";

const { Header, Footer, Sider, Content } = Layout;

const StyledPageHeader = styled(PageHeader)`
  
`

export default function Landing() {
  return (
    <Layout>
      <Header>
        Hi
      </Header>
      <Content>Content</Content>
      <Footer>Footer</Footer>
    </Layout>)
}