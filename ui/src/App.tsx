import React from 'react';

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
  Link,
  useHistory
} from "react-router-dom";

import Landing from "./views/Landing";
import Login from "./views/Login";
import Register from "./views/Register";
import './App.css';

import jwt_decode from 'jwt-decode'
import AuthContext, { AuthContextValues } from './contexts/auth'
import { AuthClaims } from './api/auth'

const { Header, Footer, Sider, Content } = Layout;
const { Title, Paragraph, Text } = Typography;

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

  let history = useHistory()
  



  return (
    <AuthContext.Provider value={auth}>
      <Router>
        <Layout style={{ minHeight: '100vh' }}>
          <PageHeader
            ghost={false}
            title={
              <Link to="/" key="logo-home">
                <a>AstraTutor</a>
              </Link>
            }
            extra={[
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
              </Link>,
            ]}
          />
          <Content>
            <Switch>
              <Route path="/" exact={true}>
                <Landing />
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
