import React, { useContext } from 'react';

import 'antd/dist/antd.css';
import { Layout, PageHeader, Button, Divider, Row, Col, Typography } from 'antd';

import { UserOutlined } from '@ant-design/icons';

import { Switch, Route, Link, useHistory, useLocation } from 'react-router-dom';

import Account from './views/Account';
import Landing from './views/Landing';
import Login from './views/Login';
import Register from './views/Register';
import Subjects from './views/Subjects';
import './App.css';

import config from './config';
import { fetchRest } from './api/rest';
import { AuthContext, useAuthValues, PrivateRoute } from './api/auth';
import { AuthClaims } from './api/auth';
import { useAsync } from 'react-async-hook';

const { Header, Footer, Sider, Content } = Layout;
const { Title, Paragraph, Text } = Typography;

function App() {
  const history = useHistory();
  const auth = useAuthValues();
  const location = useLocation();

  React.useEffect(() => {
    try {
      console.log('hi');
      auth.loginSilent();
    } catch (e) {
      console.error(`error attempting to login from localStorage ${e}`);
    }
  }, []);

  // If login state has changed, check if their profile exists
  useAsync(async () => {
    if (auth.isLoggedIn()) {
      // Check if the user has a profile
      try {
        const res = await fetchRest(
          `${config.apiUrl}/${auth.account.type}s/${auth.claims.sub}/profile`,
          {
            headers: {
              Authorization: `Bearer ${auth.bearerToken}`,
            },
          },
          [200, 404],
        );

        // No profile, redirect   them
        if (res.status === 404 && location.pathname !== '/account/profile/create') {
          history.replace('/account/profile/create');
        }
      } catch (e) {
        console.log(e);
        // TODO(ocanty) - errorhandling
      }
    }
  }, [auth, location.pathname]);

  // Don't render the page until the silent login attempt is finished
  if (!auth.loginSilentFinished()) return <AuthContext.Provider value={auth}></AuthContext.Provider>;

  let headerLinks = [];
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
      <Link to="/account/profile" key="account">
        <Button type="primary">
          <UserOutlined />
          Account
        </Button>
      </Link>,
      <Button key="logout" onClick={() => auth.logout()}>
        Logout
      </Button>,
    ];
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
      </Link>,
    ];
  }

  return (
    <AuthContext.Provider value={auth}>
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
            <PrivateRoute path="/account" component={Account} />

            <Route exact path="/subjects">
              <Subjects />
            </Route>
            <Route path="/subjects/:subject_slug/tutors"></Route>
            <Route path="/tutors/:slug"></Route>
            <Route path="/tutors/:slug/profile"></Route>
            <PrivateRoute path="/lessons" />
            <PrivateRoute path="/lessons/:lid" />
            <PrivateRoute path="/lessons/:lid/lobby" />
            <PrivateRoute path="/lessons/:lid/classroom" />
            <Route path="/login" component={Login} />
            <Route path="/register" component={Register} />
          </Switch>
        </Content>
        <Footer>
          <Divider orientation="left">AstraTutor</Divider>
          <Row>
            <Col flex={16}>Site Map</Col>
            <Col flex={24 - 16}>Links</Col>
          </Row>
          <Row style={{ margin: '0 auto', textAlign: 'center' }}>
            <p>Made with love by CS3505 Team 4</p>
          </Row>
        </Footer>
      </Layout>
    </AuthContext.Provider>
  );
}

export default App;
