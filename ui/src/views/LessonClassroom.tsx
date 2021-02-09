import React, { ReactElement, useContext } from 'react';
import styled from 'styled-components';
import { Button, Col, Layout, Modal, Row, Select, Typography } from 'antd';
import { SettingsCTX } from '../services/classroom';
import { useAsync } from 'react-async-hook';
import { ReadProfileDTO } from '../api/definitions';
import { AuthContext } from '../api/auth';
import { GetProfile } from '../services/profile';
import { UserAvatar } from '../components/UserAvatar';
import Messaging from '../components/Messaging';
import { CameraFilled, PhoneFilled, SettingFilled } from '@ant-design/icons';

interface IWebcam {
  profile: ReadProfileDTO;
  ref: JSX.Element;
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

export default function LessonClassroom(): ReactElement {
  const settings = useContext(SettingsCTX);
  const auth = useContext(AuthContext);
  const [webcamDisplays, setWebcamDisplays] = React.useState<IWebcam[]>([]);
  const [settingsOpen, setSettingsOpen] = React.useState(false);

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
  }, [settings.selectedWebcam]);
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
                <PhoneFilled />
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
          <Messaging height={window.innerHeight - webcamDisplays.length * webcamHeight} />
        </StyledSider>
        <Layout.Content></Layout.Content>
      </StyledLayout>
    </StyledLayout>
  );
}
