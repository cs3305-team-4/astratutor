import React from 'react';
import styled from 'styled-components';

import 'antd/dist/antd.css';
import {
  Layout,
  PageHeader,
  Button,
  Divider,
  Row,
  Col,
  Typography
} from "antd";

import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";

import Landing from "./views/Landing"
import './App.css';

import jwt_decode from 'jwt-decode'
import AuthContext, { AuthContextValues } from './contexts/auth'
import { AuthClaims } from './api/auth'

const { Header, Footer, Sider, Content } = Layout;
const { Title, Paragraph, Text } = Typography;

const StyledLayoutHeader = styled(Header)`
  background: white;
`

const StyledPageHeader = styled(PageHeader)``

function App() {
  const [auth, setAuth] = React.useState<AuthContextValues>()

  React.useEffect(() => {
  
    const loginFromJwt = (jwt: string) => {
      try {
        const claims = jwt_decode(jwt) as AuthClaims;

        setAuth({
          claims: claims,
          isLoggedIn: () => true,
          loginFromJwt: loginFromJwt
        })
      } catch (e) {
        console.error("tried to decode jwt from login but it didn't work: ", e)
      }
    }

    // TODO - load token from localstorage
    setAuth({
      claims: undefined,
      isLoggedIn: () => false,
      loginFromJwt: loginFromJwt
    })
  }, [])

  return (
    <AuthContext.Provider value={auth}>
      <Layout>
        <StyledLayoutHeader>
          <StyledPageHeader
            ghost={false}
            onBack={() => window.history.back()}
            title="AstraTutor"
            extra={[
              <Button key="2">Login</Button>,
              <Button key="1" type="primary">
                Register
              </Button>
            ]}
          >
          </StyledPageHeader>
        </StyledLayoutHeader>
        <Content>
          <Router>
            <Switch>
              <Route path="/">
                <Landing></Landing>
              </Route>
              <Route path="/account">
              </Route>
              <Route path="/subjects">
              </Route>
              <Route path="/subjects/:subject_slug/tutors">
              </Route>
              <Route path="/tutors/:slug">
              </Route>
              <Route path="/tutors/:slug/profile">
              </Route>
              <Route path="/lessons">
              </Route>
              <Route path="/lessons/:lid">
              </Route>
              <Route path="/lessons/:lid/lobby">
              </Route>
              <Route path="/lessons/:lid/classroom">
              </Route>

              <Route path="/login">
              </Route>
              <Route path="/register">
              </Route>
            </Switch>
          </Router>
        </Content>
        <Footer>
          <Divider orientation="left">AstraTutor</Divider>
          <Row>
            <Col flex={16}>
              Site Map
            </Col>
            <Col flex={24-16}>
              Links
            </Col>
          </Row>
          <Row style={{margin: "0 auto", textAlign: "center"}}>
            <p>Made with love by CS3505 Team 4</p>
          </Row>
        </Footer>
      </Layout>
    </AuthContext.Provider>
  )
}

export default App;
