import React, { ReactElement, useContext, useState } from 'react';

import 'antd/dist/antd.css';
import { Layout, PageHeader, Button, Divider, Row, Col } from 'antd';

import { UserOutlined } from '@ant-design/icons';

import { Switch, Route, Link, useHistory, useLocation } from 'react-router-dom';

import { Elements } from '@stripe/react-stripe-js';
import { loadStripe } from '@stripe/stripe-js';

import { Account } from './views/Account';
import { Landing } from './views/Landing';
import { Login } from './views/Login';
import { Subjects } from './views/Subjects';
import { Lessons } from './views/Lessons';
import { ViewProfile } from './views/ViewProfile';
import { Register } from './views/Register';
import { Tutors } from './views/Tutors';
import { LessonLobby } from './views/LessonLobby';
import './App.css';

import { APIContext, useApiValues, PrivateRoute } from './api/api';

import config from './config';

import { useAsync } from 'react-async-hook';
import { Profile } from './components/Profile';
import { ProfileResponseDTO } from './api/definitions';
import { UserAvatar } from './components/UserAvatar';
import styled from 'styled-components';

const stripePromise = loadStripe(config.stripePublishableKey);

const { Footer, Content } = Layout;

const StyledLogo = styled.img`
  height: 40px;
  position: absolute;
  object-fit: cover;
  width: 330px;
  top: 20px;
  left: -80px;
`;

function App(): React.ReactElement {
  const history = useHistory();
  const api = useApiValues();
  const location = useLocation();
  const [profile, setProfile] = useState<ProfileResponseDTO>();

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
      setProfile(await api.services.readProfileByAccountID(api.account.id, api.account?.type));
    }
  }, [api, location.pathname]);

  // Don't render the page until the silent login attempt is finished
  if (!api.loginSilentFinished()) return <APIContext.Provider value={api}></APIContext.Provider>;

  let headerLinks: ReactElement[] = [];
  if (api.isLoggedIn() && profile) {
    headerLinks = [
      <Link to="/" key="home">
        <Button type={history.location.pathname === '/' ? 'link' : 'text'}>Home</Button>
      </Link>,
      <Link to="/subjects" key="subjects">
        <Button type={history.location.pathname === '/subjects' ? 'link' : 'text'}>Subjects</Button>
      </Link>,
      <Link to="/subjects/tutors" key="tutors">
        <Button type={history.location.pathname.startsWith('/subjects/tutors') ? 'link' : 'text'}>Find A Tutor</Button>
      </Link>,
      <Link to="/lessons" key="lessons">
        <Button type={history.location.pathname.startsWith('/lessons') ? 'link' : 'text'}>My Lessons</Button>
      </Link>,
      <Link to="/account/profile" key="account">
        <Divider type="vertical" style={{ borderLeft: '1px solid rgb(169 169 169)', marginRight: 20 }} />
        <Button type={history.location.pathname.startsWith('/account') ? 'link' : 'text'}>
          <UserAvatar props={{ size: 20, style: { marginRight: 7 } }} profile={profile} />
          {profile?.first_name} {profile?.last_name}
        </Button>
      </Link>,
      <Button key="logout" onClick={() => api.logout()}>
        Logout
      </Button>,
    ];
  } else {
    headerLinks = [
      <Link to="/" key="home">
        <Button type={history.location.pathname === '/' ? 'link' : 'text'}>Home</Button>
      </Link>,
      <Link to="/subjects" key="subjects">
        <Button type={history.location.pathname === '/subjects' ? 'link' : 'text'}>Subjects</Button>
      </Link>,
      <Link to="/subjects/tutors" key="tutors">
        <Button type={history.location.pathname.startsWith('/subjects/tutors') ? 'link' : 'text'}>Find A Tutor</Button>
      </Link>,
      <Link to="/login" key="login">
        <Divider type="vertical" style={{ borderLeft: '1px solid rgb(169 169 169)', marginRight: 20 }} />
        <Button style={{ marginLeft: 20 }} type="primary">
          Log in
        </Button>
      </Link>,
      <Link to="/register" key="register">
        <Button>Register</Button>
      </Link>,
    ];
  }

  return (
    <APIContext.Provider value={api}>
      <APIContext.Provider value={api}>
        <Elements stripe={stripePromise}>
          <Layout style={{ minHeight: '100vh' }}>
            <PageHeader
              ghost={false}
              title={
                <Link to="/" key="logo-home">
                  <span>
                    <StyledLogo title="AstraTutor" src="/logo.svg" alt="AstraTutor" />
                  </span>
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
                <Route exact path="/subjects">
                  <Subjects />
                </Route>
                <Route path="/subjects/tutors">
                  <Tutors />
                </Route>
                <Route path="/subjects"></Route>
                <Route path="/subjects/:subject_slug/tutors"></Route>
                <Route exact path="/tutors/:slug"></Route>
                <Route exact path="/tutors/:uuid/profile" component={ViewProfile} />
                <PrivateRoute exact path="/lessons" component={Lessons} />
                <PrivateRoute exact path="/lessons/:lid" component={Lessons} />
                <PrivateRoute
                  exact
                  path={['/lessons/:lid/lobby', '/lessons/:lid/classroom', '/lessons/:lid/goodbye']}
                  component={LessonLobby}
                />
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
        </Elements>
      </APIContext.Provider>
    </APIContext.Provider>
  );
}

export default App;
