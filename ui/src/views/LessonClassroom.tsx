import {
  AudioOutlined,
  CameraFilled,
  DesktopOutlined,
  PhoneFilled,
  SettingFilled,
  VideoCameraOutlined,
} from '@ant-design/icons';
import { Button, Col, Layout, Modal, Row, Select, Tooltip, Typography } from 'antd';
import React, { ReactElement, useContext } from 'react';
import { useAsync } from 'react-async-hook';
import { useHistory, useParams } from 'react-router-dom';
import styled from 'styled-components';
import { AuthContext } from '../api/auth';
import { ReadProfileDTO } from '../api/definitions';
import Messaging from '../components/Messaging';
import { UserAvatar } from '../components/UserAvatar';
import { SettingsCTX } from '../services/classroom';
import { GetProfile } from '../services/profile';

interface IWebcam {
  profile: ReadProfileDTO;
  ref: React.ReactElement<HTMLVideoElement>;
}

const webcamHeight = 200;

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
  padding: 0;
  margin: 0;
  overflow: hidden;
  height: ${webcamHeight}px;
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
    position: absolute;
    width: 100%;
    top: 158px;
  }
`;

const StyledIcon = styled(SettingFilled)`
  color: rgb(192, 192, 192);
  font-size: 30px;
  position: fixed;
  top: 70px;
  right: 15px;
  cursor: pointer;
`;

const StyledTools = styled(Layout.Footer)`
  background-color: rgb(5 5 5);
  position: absolute;
  bottom: 0;
  left: 300px;
  width: calc(100% - 300px);
  text-align: center;
  vertical-align: middle;
`;

export default function LessonClassroom(): ReactElement {
  const { lid } = useParams<{ lid: string }>();
  const settings = useContext(SettingsCTX);
  const history = useHistory();
  const auth = useContext(AuthContext);
  const [webcamDisplays, setWebcamDisplays] = React.useState<IWebcam[]>([]);
  const [settingsOpen, setSettingsOpen] = React.useState(false);
  const [webcamEnabled, setWebcamEnabled] = React.useState(true);
  const [screenEnabled, setScreenEnabled] = React.useState(false);
  const [micEnabled, setMicEnabled] = React.useState(true);

  const addWebcam = (web: IWebcam) => {
    const other = webcamDisplays.findIndex((v) => v.ref.key === web.ref.key);
    if (other > -1) {
      const temp = webcamDisplays;
      delete temp[other];
      temp[other] = web;
      setWebcamDisplays(temp);
    } else {
      setWebcamDisplays(webcamDisplays.concat(web));
    }
  };

  useAsync(async () => {
    // Get Self Profile
    if (!webcamEnabled) {
      settings.webcamStream?.getVideoTracks().forEach((v) => {
        v.enabled = false;
      });
    } else {
      settings.webcamStream?.getVideoTracks().forEach((v) => {
        v.enabled = true;
      });
    }
    const video = (
      <video
        key={'self'}
        ref={(ref) => {
          if (ref) {
            ref.srcObject = webcamEnabled ? settings.webcamStream : null;
            webcamEnabled && ref.play();
            navigator.mediaDevices.getSupportedConstraints();
          }
        }}
      />
    );
    const web: IWebcam = { profile: await GetProfile(auth), ref: video };
    addWebcam(web);
  }, [settings.webcamStream, webcamEnabled]);

  const hangup = () => {
    settings.webcamStream?.getVideoTracks().forEach((v) => {
      v.stop();
    });
    history.push(`/lessons/${lid}/goodbye`);
  };

  return (
    <StyledLayout>
      <StyledLayout>
        <StyledIcon onClick={() => setSettingsOpen(true)} title="Settings" />
        <Modal
          title="Settings"
          visible={settingsOpen}
          centered
          onOk={() => setSettingsOpen(false)}
          onCancel={() => setSettingsOpen(false)}
        >
          <Layout>
            <br />
            <Row align="middle" justify="center">
              <Col>
                <CameraFilled />
                <Select
                  style={{ width: 300, marginLeft: 10 }}
                  value={
                    settings.selectedWebcam || (settings.webcams.length ? settings.webcams[0].deviceId : undefined)
                  }
                  onSelect={(id) => {
                    settings.setSelectedWebcam(id as string);
                  }}
                  placeholder="Select a Camera"
                >
                  {(() => {
                    const opts: ReactElement[] = [];
                    for (const dev of settings.webcams) {
                      opts.push(<Select.Option value={dev.deviceId}>{dev.label}</Select.Option>);
                    }
                    return opts;
                  })()}
                </Select>
              </Col>
            </Row>
            <br />
            <Row align="middle" justify="center">
              <Col>
                <AudioOutlined />
                <Select
                  style={{ width: 300, marginLeft: 10 }}
                  value={
                    settings.selectedMicrophone ||
                    (settings.microphones.length ? settings.microphones[0].deviceId : undefined)
                  }
                  onSelect={(id) => {
                    settings.setSelectedMicrophone(id as string);
                  }}
                  placeholder="Select a Microphone"
                >
                  {(() => {
                    const opts: ReactElement[] = [];
                    for (const dev of settings.microphones) {
                      opts.push(<Select.Option value={dev.deviceId}>{dev.label}</Select.Option>);
                    }
                    return opts;
                  })()}
                </Select>
              </Col>
            </Row>
            <br />
          </Layout>
        </Modal>
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
          <Messaging height={webcamDisplays.length * webcamHeight} />
        </StyledSider>
        <Layout.Content></Layout.Content>
        <StyledTools>
          <Tooltip title="Toggle Mute">
            <Button
              ghost={!micEnabled}
              onClick={() => setMicEnabled(!micEnabled)}
              size={'large'}
              shape="circle"
              style={{ margin: '0 10px' }}
            >
              <AudioOutlined size={20} style={{ color: micEnabled ? '#000' : '#fff' }} />
            </Button>
          </Tooltip>
          <Tooltip title="Toggle Webcam">
            <Button
              ghost={!webcamEnabled}
              onClick={() => setWebcamEnabled(!webcamEnabled)}
              size={'large'}
              shape="circle"
              style={{ margin: '0 10px' }}
            >
              <VideoCameraOutlined size={20} style={{ color: webcamEnabled ? '#000' : '#fff' }} />
            </Button>
          </Tooltip>
          <Tooltip title="Share Screen">
            <Button
              ghost={!screenEnabled}
              onClick={() => setScreenEnabled(!screenEnabled)}
              size={'large'}
              shape="circle"
              style={{ margin: '0 10px' }}
            >
              <DesktopOutlined size={20} style={{ color: screenEnabled ? '#000' : '#fff' }} />
            </Button>
          </Tooltip>
          <Tooltip title="Hang Up">
            <Button
              onClick={hangup}
              size={'large'}
              shape="circle"
              style={{ backgroundColor: '#c50505', margin: '0 10px' }}
            >
              <PhoneFilled size={20} rotate={225} style={{ color: '#fff' }} />
            </Button>
          </Tooltip>
        </StyledTools>
      </StyledLayout>
    </StyledLayout>
  );
}
