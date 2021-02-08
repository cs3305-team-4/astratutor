import { Layout, Button, Typography, Avatar, Tooltip, Col, Row, Divider } from 'antd';
import React, { ReactElement } from 'react';
import { RouteComponentProps, useLocation, useParams } from 'react-router-dom';
import styled from 'styled-components';

const StyledLayout = styled(Layout)`
  background-color: rgb(21 20 20);
  padding: 30vw;
  color: #fff;
`;

const StyledDivider = styled(Divider)`
  border-top: 1px solid rgb(255 252 252 / 11%);
`;

export default function LessonLobby(): ReactElement {
  const { lid } = useParams<{ lid: string }>();
  const location = useLocation();
  return (
    <StyledLayout>
      <Typography>
        <Typography.Title style={{ color: '#fff', textAlign: 'center' }} level={1}>
          Joining your Mathematics 101 lesson!
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
      <StyledDivider />
      <Button style={{ width: '50%', margin: 'auto' }} ghost type="primary">
        Join
      </Button>
    </StyledLayout>
  );
}
