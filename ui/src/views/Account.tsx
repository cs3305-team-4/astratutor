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
import CreateProfileForm from '../components/CreateProfileForm'
import { AuthContext } from '../api/auth';


const { SubMenu } = Menu;
const { Title, Paragraph, Text } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const StyledLayout = styled(Layout)`
  background-color: rgba(233,233,233);
`

const StyledSider = styled(Sider)`
  background-color: rgba(233,233,233);
`

export default function Account() {
  const history = useHistory()
  let { path, url } = useRouteMatch()

  const location = useLocation()

  return (
    <StyledLayout>
      <StyledSider>
        <Menu
          selectedKeys={[location.pathname.substr(path.length)]}
          mode="inline"
        >
          <Menu.Item onClick={()=>history.push(`${path}/profile`)} key="/profile" icon={<UserOutlined />}>
            Profile
          </Menu.Item>
          <Menu.Item onClick={()=>history.push(`${path}/billing`)} key="/billing" icon={<BankOutlined />}>
            Billing
          </Menu.Item>
        </Menu>
      </StyledSider>
      <Content style={{minHeight: "calc(100vh - 72px)"}}>
        <Switch>
          <Route exact path={`${path}/profile/create`}>
            <Row>
              <Col md={9} sm={6} xs={0}/>
              <Col md={24} sm={24} xs={24} style={{padding: "0.5rem", backgroundColor: "rgba(255,255,255,0.8)"}}>
                <AuthContext.Consumer>
                  {auth => auth.isLoggedIn() ? <CreateProfileForm uuid={auth.claims.sub} type={auth.account.type}/> : <></>}
                </AuthContext.Consumer>
              </Col>
              <Col md={9} sm={6} xs={0}/>
            </Row>
          </Route>
          <Route exact path={`${path}/profile`}>
            <Row>
              <Col md={24} sm={24} xs={24} style={{padding: "0.5rem", backgroundColor: "rgba(255,255,255,0.8)"}}>
                <AuthContext.Consumer>
                  {auth => auth.isLoggedIn() ? <AccountProfile uuid={auth.claims.sub} type={auth.account.type}/> : <></>}
                </AuthContext.Consumer>
              </Col>
            </Row>
          </Route>
          <Route exact path={`${path}/billing`}>
            <Row>
              <Col md={24} sm={10} xs={24} style={{padding: "1rem", backgroundColor: "rgba(255,255,255,0.8)"}}>
                Billing
              </Col>
            </Row>
          </Route>
        </Switch>
      </Content>
    </StyledLayout>
  )
}