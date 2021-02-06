import React, {useContext} from 'react';

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
  UserOutlined
} from "@ant-design/icons"

import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link,
  useHistory
} from "react-router-dom";

import Landing from "./views/Landing";
import Login from "./views/Login";
import Register from "./views/Register";
import './App.css';

import { AuthContext, useAuthValues, PrivateRoute } from './api/auth'
import { AuthClaims } from './api/auth'

const { Header, Footer, Sider, Content } = Layout;
const { Title, Paragraph, Text } = Typography;

function App() {
  let auth = useAuthValues()

  React.useEffect(() => {
    try {
      auth.loginFromLocalStorage()
    } catch (e) {
      console.error(`error attempting to login from localStorage ${e}`)
    }
  }, [])

  let headerLinks = []
  if (auth.isLoggedIn()) {
    headerLinks = [
      <Link to="/" key="home">
        <Button type="text">Home</Button>
      </Link>,
      <Link to="/subjects" key="subjects">
        <Button type="text">Subjects</Button>
      </Link>,
      <Link to="/subjects/tutors" key="tutors">
        <Button type="text">Find A Tutor</Button>
      </Link>,
      <Link to="/lessons" key="lessons">
        <Button type="text">My Lessons</Button>
      </Link>,
      <Link to="/account" key="profile">
        <Button type="primary">
          <UserOutlined />
          Account
        </Button>
      </Link>,
      <Button onClick={()=>auth.logout()}>Logout</Button>
    ]
  } else {
    headerLinks = [
      <Link to="/" key="home">
        <Button type="text">Home</Button>,
      </Link>,
      <Link to="/subjects" key="subjects">
        <Button type="text">Subjects</Button>,
      </Link>,
      <Link to="/subjects/tutors" key="tutors">
        <Button type="text">Find A Tutor</Button>,
      </Link>,
      <Link to="/login" key="login">
        <Button type="primary">Log in</Button>
      </Link>,
      <Link to="/register" key="register">
        <Button>Register</Button>
      </Link>
    ]
  }

  return (
    <AuthContext.Provider value={auth}>
      <Router>
        <Layout style={{ minHeight: '100vh' }}>
          <PageHeader
            ghost={false}
            title={
              <Link to="/" key="logo-home">
                <span>AstraTutor</span>
              </Link>
            }
            extra={headerLinks}
          />
          <Content>
            <Switch>
              <Route path="/" exact={true}>
                <Landing />
              </Route>
              <PrivateRoute path="/account"/>
              <PrivateRoute path="/account/profile"/>
              <Route path="/subjects">
              </Route>
              <Route path="/subjects/:subject_slug/tutors">
              </Route>
              <Route path="/tutors/:slug">
              </Route>
              <Route path="/tutors/:slug/profile">
              </Route>
              <PrivateRoute path="/lessons"/>
              <PrivateRoute path="/lessons/:lid"/>
              <PrivateRoute path="/lessons/:lid/lobby"/>
              <PrivateRoute path="/lessons/:lid/classroom"/>
              <Route path="/login" component={Login}/>
              <Route path="/register" component={Register}/>
            </Switch>
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
      </Router>
    </AuthContext.Provider>
  )
}

export default App;
