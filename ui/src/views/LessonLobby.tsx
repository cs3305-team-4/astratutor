import {
  ArrowLeftOutlined,
  AudioOutlined,
  CameraFilled,
  FullscreenExitOutlined,
  FullscreenOutlined,
  PhoneFilled,
} from '@ant-design/icons';
import { Avatar, Button, Col, Divider, Layout, Row, Select, Tooltip, Typography } from 'antd';
import React, { ReactElement, useContext, useEffect, useRef, useState } from 'react';
import { useAsync } from 'react-async-hook';
import { Link, Route, Switch, useHistory, useParams } from 'react-router-dom';
import styled from 'styled-components';
import { APIContext } from '../api/api';
import { UserAvatar } from '../components/UserAvatar';
import { ISettings, SettingsCTX } from '../api/classroom';
import { LessonClassroom } from './LessonClassroom';
import { MESSAGE_TYPE, Signalling } from '../webrtc/signalling';
import * as Devices from '../webrtc/devices';
import config from '../config';
import {
  AccountType,
  LessonRequestDTO,
  LessonResponseDTO,
  ProfileResponseDTO,
  SubjectTaughtDTO,
} from '../api/definitions';

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

export function LessonLobby(): ReactElement {
  const { lid } = useParams<{ lid: string }>();
  const api = useContext(APIContext);
  const history = useHistory();
  const signalling = useRef<Signalling>();
  const [webcams, setWebcams] = useState<MediaDeviceInfo[]>([]);
  const [microphones, setMicrophones] = useState<MediaDeviceInfo[]>([]);
  const display = useRef<HTMLVideoElement>();
  const [selectedWebcam, setSelectedWebcam] = useState<string>('');
  const [selectedMicrophone, setSelectedMicrophone] = useState<string>('');
  const [fullscreen, setFullscreen] = useState(document.fullscreenElement !== null);
  const [joined, setJoined] = useState(false);
  const [webcamStream, setWebcamStream] = useState<MediaStream | null>(null);
  const [otherProfiles, setOtherProfiles] = React.useState<{ [id: string]: ProfileResponseDTO }>({});
  const [metadata, setMetadata] = useState<LessonResponseDTO>();
  const [completed, setCompleted] = useState(false);

  if (!joined && !history.location.pathname.endsWith('lobby')) {
    history.push(`/lessons/${lid}/lobby`);
  }

  const settingsValue: ISettings = {
    signalling: signalling.current,
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
    otherProfiles,
  };

  useEffect(() => {
    signalling.current = new Signalling(api.claims?.sub ?? '', `${config.signallingUrl}/${lid}`, {
      onopen: (event: Event) => {
        console.log('Connected to WS: ', lid);
        // TODO(james): Probe Users
        // signalling.current?.send(MESSAGE_TYPE.PROBE, '', null);
      },
      onclose: console.log,
      onerror: console.log,
    });
  }, []);

  useAsync(async () => {
    const getDevices = async () => {
      await Devices.devicePermissions();
      const devices = await Devices.getDevices();
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
      const lesson = await api.services.readLesson(lid);
      setOtherProfiles({
        [lesson.student_id]: await api.services.readProfileByAccountID(lesson.student_id, AccountType.Student),
        [lesson.tutor_id]: await api.services.readProfileByAccountID(lesson.tutor_id, AccountType.Tutor),
      });
      setMetadata(lesson);
    };
    getDevices();
    navigator.mediaDevices.ondevicechange = getDevices;
  }, []);
  useAsync(async () => {
    if (!webcamStream) {
      await Devices.devicePermissions();
      const devices = await Devices.getDevices();
      const mic = devices.filter((v) => v.kind === 'audioinput');
      const mid = mic.length ? mic[0].deviceId : '';
      console.log('MIC', mid);
      setSelectedMicrophone(mid);
      const dev = devices.filter((v) => v.kind === 'videoinput');
      const id = dev.length ? dev[0].deviceId : '';
      setSelectedWebcam(id);
      const stream = await Devices.cameraStream(id, mid);
      setWebcamStream(stream);
      console.log('webcam init');
    }
  }, []);
  useEffect(() => {
    return () => {
      webcamStream?.getTracks().forEach((v) => {
        v.stop();
      });
    };
  }, [webcamStream]);

  useAsync(async () => {
    if (webcamStream) {
      webcamStream.getVideoTracks().forEach((v) => v.stop());
      webcamStream.getAudioTracks().forEach((v) => v.stop());
      setWebcamStream(await Devices.cameraStream(selectedWebcam, selectedMicrophone));
      console.log('webcam change');
    }
  }, [selectedWebcam, selectedMicrophone]);
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
              Thanks for attending {metadata?.subject_name}!
            </Typography.Title>
            {(api.account?.type === AccountType.Student || completed) && (
              <Button style={{ width: '50%', margin: '.3em auto' }} ghost type="link">
                <Link to={`/lessons/completed?reschedule=${metadata?.id}`}>Schedule my next lesson</Link>
              </Button>
            )}
            {api.account?.type === AccountType.Student && (
              <Button style={{ width: '', margin: '.3em auto' }} type="primary">
                <Link to={`/tutors/${metadata?.tutor_id}/profile`}>Rate my tutor</Link>
              </Button>
            )}

            {api.account?.type === AccountType.Tutor && !completed && (
              <Button
                onClick={async () => {
                  await api.services.updateLessonStageCompleted(metadata?.id ?? '');
                  setCompleted(true);
                }}
                style={{ width: '', margin: '.3em auto' }}
                type="primary"
              >
                Mark lesson as completed
              </Button>
            )}

            {(api.account?.type === AccountType.Student || completed) && (
              <>
                <StyledDivider />
                <Button style={{ width: '50%', margin: '.1em auto' }} ghost type="primary">
                  <Link to={`/lessons`}>Go back to my lessons</Link>
                </Button>
              </>
            )}
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
                Joining your {metadata?.subject_name} classroom!
              </Typography.Title>
            </Typography>
            {/* TODO(james): Send probe message to discover users */}
            <Typography style={{ textAlign: 'center' }}>
              <Typography.Text style={{ color: '#fff' }}>Already in this meeting:</Typography.Text>
            </Typography>
            <Row align="middle" justify="center">
              <Col>
                <Avatar.Group size="default">
                  <Tooltip title="Gamer">
                    <UserAvatar
                      profile={{
                        account_id: '1',
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
                      opts.push(
                        <Option key={dev.deviceId} value={dev.deviceId}>
                          {dev.label}
                        </Option>,
                      );
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
                      opts.push(
                        <Option key={dev.deviceId} value={dev.deviceId}>
                          {dev.label}
                        </Option>,
                      );
                    }
                    return opts;
                  })()}
                </StyledSelect>
              </Col>
            </Row>
            <br />
            <video
              muted
              style={{
                height: 300,
              }}
              ref={async (r) => {
                if (r && webcamStream) {
                  try {
                    if (webcamStream.getTracks().filter((v) => v.enabled).length > 0) {
                      r.srcObject = webcamStream;
                      await r.play();
                    }
                  } catch ({ message }) {
                    console.error(message);
                  }
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
