import {
  ArrowLeftOutlined,
  AudioOutlined,
  CameraFilled,
  FullscreenExitOutlined,
  FullscreenOutlined,
  PhoneFilled,
} from '@ant-design/icons';
import { Avatar, Button, Col, Divider, Layout, Row, Select, Tooltip, Typography } from 'antd';
import React, { ReactElement, useEffect, useRef, useState } from 'react';
import { useAsync } from 'react-async-hook';
import { Link, Route, Switch, useHistory, useParams } from 'react-router-dom';
import styled from 'styled-components';
import { UserAvatar } from '../components/UserAvatar';
import { ISettings, SettingsCTX } from '../services/classroom';
import LessonClassroom from './LessonClassroom';

const { Option } = Select;

const StyledNav = styled.nav`
  position: fixed;
  right: 0;
  top: 5px;
  z-index: 200;
  color: #fff;
`;

const StyledLayout = styled(Layout)`
  background-color: rgb(21 20 20);
  padding: 5em 35vw;
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
const StyledDivider = styled(Divider)`
  border-top: 1px solid rgb(255 252 252 / 11%);
`;

const StyledSelect = styled(Select)`
  padding-left: 20px;
  margin: auto;
  width: 200px;

  & .ant-select-selector {
    background-color: transparent !important;
  }
  & .ant-select-arrow,
  & .ant-select-selection-item {
    color: #fff;
  }
`;

export default function LessonLobby(): ReactElement {
  const { lid } = useParams<{ lid: string }>();
  const history = useHistory();
  const [webcams, setWebcams] = useState<MediaDeviceInfo[]>([]);
  const [microphones, setMicrophones] = useState<MediaDeviceInfo[]>([]);
  const display = useRef<HTMLVideoElement>();
  const [selectedWebcam, setSelectedWebcam] = useState<string>('');
  const [selectedMicrophone, setSelectedMicrophone] = useState<string>('');
  const [fullscreen, setFullscreen] = useState(document.fullscreenElement !== null);
  const [joined, setJoined] = useState(false);
  const [webcamStream, setWebcamStream] = useState<MediaStream | null>(null);

  if (!joined && !history.location.pathname.endsWith('lobby')) {
    history.push(`/lessons/${lid}/lobby`);
  }

  const settingsValue: ISettings = {
    fullscreen,
    setFullscreen,
    webcams,
    setWebcams,
    microphones,
    setMicrophones,
    selectedWebcam,
    setSelectedWebcam,
    selectedMicrophone,
    setSelectedMicrophone,
    webcamStream,
    setWebcamStream,
  };

  useAsync(async () => {
    await navigator.mediaDevices.getUserMedia({ video: true });
    const devices = await navigator.mediaDevices.enumerateDevices();
    const vid: MediaDeviceInfo[] = [];
    const mic: MediaDeviceInfo[] = [];
    for (const dev of devices) {
      switch (dev.kind) {
        case 'videoinput':
          vid.push(dev);
          break;
        case 'audioinput':
          mic.push(dev);
          break;
      }
      setWebcams(vid);
      setMicrophones(mic);
    }
  }, []);
  useAsync(async () => {
    if (!webcamStream) {
      await navigator.mediaDevices.getUserMedia({ video: true });
      const devices = await navigator.mediaDevices.enumerateDevices();
      const dev = devices.filter((v) => v.kind === 'videoinput');
      const id = dev.length ? dev[0].deviceId : '';
      setSelectedWebcam(id);
      const stream = await navigator.mediaDevices.getUserMedia({ video: { deviceId: selectedWebcam } });
      setWebcamStream(stream);
    }
  }, []);

  useAsync(async () => {
    if (webcamStream) {
      webcamStream.getVideoTracks().forEach((v) => v.stop());
      setWebcamStream(await navigator.mediaDevices.getUserMedia({ video: { deviceId: selectedWebcam } }));
    }
  }, [selectedWebcam]);
  const [title, setTitle] = useState('Mathematics 101');
  return (
    <SettingsCTX.Provider value={settingsValue}>
      <Switch>
        <Route path="/lessons/:lid/goodbye">
          <StyledNav>
            <Button
              type="link"
              ghost
              onClick={() => {
                if (document.fullscreenElement) {
                  document.exitFullscreen();
                  setFullscreen(false);
                } else {
                  document.documentElement.requestFullscreen();
                  setFullscreen(true);
                }
              }}
            >
              {fullscreen ? (
                <FullscreenExitOutlined title="Exit Fullscreen" style={{ color: '#c0c0c0', fontSize: 30 }} />
              ) : (
                <FullscreenOutlined title="Fullscreen" style={{ color: '#c0c0c0', fontSize: 30 }} />
              )}
            </Button>
          </StyledNav>
          <StyledLayout>
            <Typography.Title style={{ color: '#fff', textAlign: 'center' }} level={1}>
              Thanks for attending {title}!
            </Typography.Title>
            <Button style={{ width: '50%', margin: '.1em auto' }} ghost type="link">
              Schedule my next lesson
            </Button>
            <StyledDivider />
            <Button style={{ width: '50%', margin: '.1em auto' }} ghost type="primary">
              <Link to={`/lessons`}>Go back to my lessons</Link>
            </Button>
          </StyledLayout>
        </Route>
        <Route path="/lessons/:lid/classroom">
          <StyledNav>
            <Button
              type="link"
              ghost
              onClick={() => {
                if (document.fullscreenElement) {
                  document.exitFullscreen();
                  setFullscreen(false);
                } else {
                  document.documentElement.requestFullscreen();
                  setFullscreen(true);
                }
              }}
            >
              {fullscreen ? (
                <FullscreenExitOutlined title="Exit Fullscreen" style={{ color: '#c0c0c0', fontSize: 30 }} />
              ) : (
                <FullscreenOutlined title="Fullscreen" style={{ color: '#c0c0c0', fontSize: 30 }} />
              )}
            </Button>
          </StyledNav>
          <LessonClassroom />
        </Route>
        <Route path="/lessons/:lid/lobby">
          <StyledNav>
            <Button
              type="link"
              ghost
              onClick={() => {
                window.history.back();
              }}
            >
              <ArrowLeftOutlined title="Go back" style={{ color: '#c0c0c0', fontSize: 30 }} />
            </Button>
            <Button
              type="link"
              ghost
              onClick={() => {
                if (document.fullscreenElement) {
                  document.exitFullscreen();
                  setFullscreen(false);
                } else {
                  document.documentElement.requestFullscreen();
                  setFullscreen(true);
                }
              }}
            >
              {fullscreen ? (
                <FullscreenExitOutlined title="Exit Fullscreen" style={{ color: '#c0c0c0', fontSize: 30 }} />
              ) : (
                <FullscreenOutlined title="Fullscreen" style={{ color: '#c0c0c0', fontSize: 30 }} />
              )}
            </Button>
          </StyledNav>
          <StyledLayout>
            <Typography>
              <Typography.Title style={{ color: '#fff', textAlign: 'center' }} level={1}>
                Joining your {title} classroom!
              </Typography.Title>
            </Typography>
            <Typography style={{ textAlign: 'center' }}>
              <Typography.Text style={{ color: '#fff' }}>Already in this meeting:</Typography.Text>
            </Typography>
            <Row align="middle" justify="center">
              <Col>
                <Avatar.Group size="default">
                  <Tooltip title="Gamer">
                    <UserAvatar
                      profile={{
                        avatar: '',
                        slug: '/',
                        first_name: 'Gamer',
                        last_name: 'Jones',
                        city: 'Cark',
                        country: 'Ireland',
                        subtitle: 'Gamer',
                        description: 'Gamer',
                        color: '#199a4c',
                      }}
                    />
                  </Tooltip>
                </Avatar.Group>
              </Col>
            </Row>
            <br />
            <Row align="middle" justify="center">
              <Col>
                <CameraFilled />
                <StyledSelect
                  value={selectedWebcam || (webcams.length ? webcams[0].deviceId : undefined)}
                  onSelect={(id) => {
                    setSelectedWebcam(id as string);
                  }}
                  placeholder="Select a Camera"
                >
                  {(() => {
                    const opts: ReactElement[] = [];
                    for (const dev of webcams) {
                      opts.push(<Option value={dev.deviceId}>{dev.label}</Option>);
                    }
                    return opts;
                  })()}
                </StyledSelect>
              </Col>
            </Row>
            <br />
            <Row align="middle" justify="center">
              <Col>
                <AudioOutlined />
                <StyledSelect
                  value={selectedMicrophone || (microphones.length ? microphones[0].deviceId : undefined)}
                  onSelect={(id) => {
                    setSelectedMicrophone(id as string);
                  }}
                  placeholder="Select a Microphone"
                >
                  {(() => {
                    const opts: ReactElement[] = [];
                    for (const dev of microphones) {
                      opts.push(<Option value={dev.deviceId}>{dev.label}</Option>);
                    }
                    return opts;
                  })()}
                </StyledSelect>
              </Col>
            </Row>
            <br />
            <video
              style={{
                height: 300,
              }}
              ref={(r) => {
                if (r && webcamStream) {
                  r.srcObject = webcamStream;
                  r.play();
                }
              }}
            ></video>
            <StyledDivider />
            <Button style={{ width: '50%', margin: '.1em auto' }} ghost type="primary">
              <Link onClick={() => setJoined(true)} to={`/lessons/${lid}/classroom`}>
                Join
              </Link>
            </Button>
          </StyledLayout>
        </Route>
      </Switch>
    </SettingsCTX.Provider>
  );
}
