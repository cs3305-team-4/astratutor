import React from 'react';

import 'antd/dist/antd.css';
import { Layout, PageHeader, Button, Divider, Row, Col } from 'antd';

import { UserOutlined } from '@ant-design/icons';

import { Switch, Route, Link, useHistory, useLocation } from 'react-router-dom';

import { Account } from './views/Account';
import { Landing } from './views/Landing';
import { Login } from './views/Login';
import { Lessons } from './views/Lessons';
import { ViewProfile } from './views/ViewProfile';
import { Register } from './views/Register';
import './App.css';

import { APIContext, useApiValues, PrivateRoute } from './api/api';

import { useAsync } from 'react-async-hook';

const { Footer, Content } = Layout;

function App(): React.ReactElement {
  const history = useHistory();
  const api = useApiValues();
  const location = useLocation();

  useAsync(async () => {
    try {
      await api.loginSilent();
    } catch (e) {
      console.error(`error attempting to login from localStorage ${e}`);
    }
  }, []);

  // If login state has changed, check if their profile exists
  useAsync(async () => {
    if (api.isLoggedIn()) {
      // Check if the user has a profile
      if (!(await api.services.accountHasProfile(api.account.id, api.account.type))) {
        if (location.pathname !== '/account/profile/create') {
          history.replace('/account/profile/create');
        }
      }
    }
  }, [api, location.pathname]);

  // Don't render the page until the silent login attempt is finished
  if (!api.loginSilentFinished()) return <APIContext.Provider value={api}></APIContext.Provider>;

  let headerLinks = [];
  if (api.isLoggedIn()) {
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
      <Button key="logout" onClick={() => api.logout()}>
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
    <APIContext.Provider value={api}>
      <APIContext.Provider value={api}>
        <Layout style={{ minHeight: '100vh' }}>
          <PageHeader
            ghost={false}
            title={
              <Link to="/" key="logo-home">
                <span>AstraTutor</span>
              </Link>
            }
            extra={headerLinks}
            style={{ boxShadow: '0 1px 10px rgba(0,0,0,0.25)' }}
          />
          <Content>
            <Switch>
              <Route path="/" exact={true}>
                <Landing />
              </Route>
              <PrivateRoute path="/account" component={Account} />
              <Route path="/subjects"></Route>
              <Route path="/subjects/:subject_slug/tutors"></Route>
              <Route exact path="/tutors/:slug"></Route>
              <Route exact path="/tutors/:uuid/profile" component={ViewProfile} />
              <PrivateRoute path="/lessons" component={Lessons} />
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
      </APIContext.Provider>
    </APIContext.Provider>
  );
}

export default App;
