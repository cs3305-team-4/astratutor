import { ReadProfileDTO } from '../api/definitions';
import { Avatar, AvatarProps, Tooltip } from 'antd';
import React from 'react';

export function UserAvatar(props: { profile: ReadProfileDTO; props?: AvatarProps }): JSX.Element {
  const backgroundColor = '#' + Math.floor(Math.random() * 16777215).toString(16);
  return (
    <Tooltip title={props.profile.first_name + ' ' + props.profile.last_name} placement="top">
      <Avatar style={{ backgroundColor, cursor: 'pointer' }} {...props.props}>
        {props.profile.avatar || (props.profile.first_name ? props.profile.first_name[0].toUpperCase() : '')}
      </Avatar>
    </Tooltip>
  );
}
