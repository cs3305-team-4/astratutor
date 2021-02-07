import React from 'react';
import styled from 'styled-components';

import {
  Link, useHistory, Switch, Route, useRouteMatch, useLocation
} from "react-router-dom"

import {
  Typography,
  Layout,
  Menu,
  Row,
  Col
} from "antd";

import {
    LockOutlined,
    MailOutlined,
    UserOutlined,
    BankOutlined,
    AppstoreOutlined,
    SettingOutlined
} from "@ant-design/icons";

import AccountProfile from '../components/AccountProfile'


const { SubMenu } = Menu;
const { Title, Paragraph, Text } = Typography;
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
  const history = useHistory()
  let { path, url } = useRouteMatch()

  const location = useLocation()

  console.log(path, location, location.pathname.substr(path.length))


  const goto = (menu: string) => {

  }


  return (
    <Layout>
      <Sider>
        <Menu
          selectedKeys={[location.pathname.substr(path.length)]}
          mode="inline"
        >
          <Menu.Item onClick={()=>history.push(`${path}/general`)} key="/general" icon={<SettingOutlined />}>
            General
          </Menu.Item>
          <Menu.Item onClick={()=>history.push(`${path}/profile`)} key="/profile" icon={<UserOutlined />}>
            Profile
          </Menu.Item>
          <Menu.Item onClick={()=>history.push(`${path}/billing`)} key="/billing" icon={<BankOutlined />}>
            Billing
          </Menu.Item>
        </Menu>
      </Sider>
      <Content>
        <Switch>
          <Route exact path={`${path}/general`}>
            <Col md={9} sm={6} xs={0}/>
            <Col md={24} sm={10} xs={24} style={{padding: "1rem", backgroundColor: "rgba(255,255,255,0.8)"}}>
              General
            </Col>
            <Col md={9} sm={6} xs={0}/>
          </Route>
          <Route path={`${path}/profile`}>
            <Row>
              <Col md={9} sm={6} xs={0}/>
              <Col md={24} sm={24} xs={24} style={{padding: "0.5rem", backgroundColor: "rgba(255,255,255,0.8)"}}>
                <AccountProfile/>
              </Col>
              <Col md={9} sm={6} xs={0}/>
            </Row>
          </Route>
          <Route path={`${path}/billing`}>
            <Row>
              <Col md={9} sm={6} xs={0}/>
              <Col md={24} sm={10} xs={24} style={{padding: "1rem", backgroundColor: "rgba(255,255,255,0.8)"}}>
                Billing
              </Col>
              <Col md={9} sm={6} xs={0}/>
            </Row>
          </Route>
        </Switch>
      </Content>
    </Layout>
  )
}