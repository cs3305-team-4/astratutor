import React, { ReactElement, useContext } from 'react';
import styled from 'styled-components';
import { Button, Layout, Typography } from 'antd';
import { SettingsCTX } from '../services/classroom';
import { FullscreenExitOutlined, FullscreenOutlined, StepBackwardOutlined } from '@ant-design/icons';
import { useAsync } from 'react-async-hook';
import { ReadProfileDTO } from '../api/definitions';
import { AuthContext } from '../api/auth';
import { GetProfile } from '../services/profile';
import { UserAvatar } from '../components/UserAvatar';

interface IWebcam {
  profile: ReadProfileDTO;
  ref: JSX.Element;
}

const StyledNav = styled.nav`
  position: fixed;
  right: 0;
  top: 5px;
  z-index: 200;
  color: #fff;
`;

const StyledLayout = styled(Layout)`
  background-color: rgb(21 20 20);
  color: #fff;
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 100;
  display: flex;
  justify-content: center;
`;

const StyledSider = styled(Layout.Sider)`
  background-color: rgb(15 15 15);
  border-right: 2px solid rgb(10 10 10);
`;

const StyledWebcam = styled.div`
  width: 100%;
  & .ant-typography {
    padding-left: 10px;
    color: #fff;
  }
  & video {
    width: 100%;
  }
  & .profile {
    background: linear-gradient(180deg, rgba(255, 255, 255, 0) 0%, rgba(0, 0, 0, 1) 100%);
    padding: 5px 10px;
    height: 42px;
    position: relative;
    top: -47px;
  }
`;

export default function LessonClassroom(): ReactElement {
  const settings = useContext(SettingsCTX);
  const auth = useContext(AuthContext);
  const [webcamDisplays, setWebcamDisplays] = React.useState<IWebcam[]>([]);

  const addWebcam = (web: IWebcam) => {
    const other = webcamDisplays.find((v) => v.ref.key === web.ref.key);
    if (!other) {
      setWebcamDisplays(webcamDisplays.concat(web));
    }
  };

  useAsync(async () => {
    // Get Self Profile
    const stream = await navigator.mediaDevices.getUserMedia({ video: { deviceId: settings.selectedWebcam } });
    const video = (
      <video
        key={'self'}
        ref={(ref) => {
          if (ref) {
            ref.srcObject = stream;
            ref.play();
            navigator.mediaDevices.getSupportedConstraints();
          }
        }}
      />
    );
    const web: IWebcam = { profile: await GetProfile(auth), ref: video };
    addWebcam(web);
  }, []);
  return (
    <StyledLayout>
      <StyledNav>
        <Button
          type="link"
          ghost
          onClick={() => {
            window.history.back();
          }}
        >
          <StepBackwardOutlined title="Go back" style={{ color: '#c0c0c0', fontSize: 30 }} />
        </Button>
        <Button
          type="link"
          ghost
          onClick={() => {
            if (document.fullscreenElement) {
              document.exitFullscreen();
              settings.setFullscreen(false);
            } else {
              document.documentElement.requestFullscreen();
              settings.setFullscreen(true);
            }
          }}
        >
          {settings.fullscreen ? (
            <FullscreenExitOutlined title="Exit Fullscreen" style={{ color: '#c0c0c0', fontSize: 30 }} />
          ) : (
            <FullscreenOutlined title="Fullscreen" style={{ color: '#c0c0c0', fontSize: 30 }} />
          )}
        </Button>
      </StyledNav>
      <StyledLayout>
        <StyledSider width={300}>
          {webcamDisplays.map((v) => (
            <StyledWebcam key={v.ref.key}>
              {v.ref}
              <div className="profile">
                <UserAvatar profile={v.profile} />
                <Typography.Text>
                  {v.profile.first_name} {v.profile.last_name}
                </Typography.Text>
              </div>
            </StyledWebcam>
          ))}
        </StyledSider>
        <Layout.Content></Layout.Content>
      </StyledLayout>
    </StyledLayout>
  );
}
