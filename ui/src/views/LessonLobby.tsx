import { CameraFilled, CameraOutlined, PhoneFilled } from '@ant-design/icons';
import { Layout, Button, Typography, Avatar, Tooltip, Col, Row, Divider, Select } from 'antd';
import React, { ReactElement, useState } from 'react';
import { useAsync } from 'react-async-hook';
import { RouteComponentProps, useLocation, useParams } from 'react-router-dom';
import styled from 'styled-components';

const { Option } = Select;

const StyledLayout = styled(Layout)`
  background-color: rgb(21 20 20);
  padding: 30vh 30vw;
  color: #fff;
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
  const location = useLocation();
  const [webcams, setWebcams] = useState<MediaDeviceInfo[]>([]);
  const [microphones, setMicrophones] = useState<MediaDeviceInfo[]>([]);
  useAsync(async () => {
    await navigator.mediaDevices.getUserMedia({ audio: true, video: true });
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
  const [title, setTitle] = useState('Mathematics 101');
  return (
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
              <Avatar style={{ backgroundColor: '#f56a00' }}>G</Avatar>
            </Tooltip>
          </Avatar.Group>
        </Col>
      </Row>
      <br />
      <Row align="middle" justify="center">
        <Col>
          <CameraFilled />
          <StyledSelect placeholder="Select a Camera">
            {(() => {
              const opts: ReactElement[] = [];
              for (const dev of webcams) {
                console.log(dev);
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
          <PhoneFilled />
          <StyledSelect placeholder="Select a Microphone">
            {(() => {
              const opts: ReactElement[] = [];
              for (const dev of microphones) {
                console.log(dev);
                opts.push(<Option value={dev.deviceId}>{dev.label}</Option>);
              }
              return opts;
            })()}
          </StyledSelect>
        </Col>
      </Row>
      <StyledDivider />
      <Button style={{ width: '50%', margin: 'auto' }} ghost type="primary">
        Join
      </Button>
    </StyledLayout>
  );
}
