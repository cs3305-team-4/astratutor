import { ReadProfileDTO } from '../api/definitions';
import { Avatar, AvatarProps } from 'antd';
import React from 'react';

export function UserAvatar(props: { profile: ReadProfileDTO; props?: AvatarProps }): JSX.Element {
  const backgroundColor = '#' + Math.floor(Math.random() * 16777215).toString(16);
  return (
    <Avatar style={{ backgroundColor }} {...props.props}>
      {props.profile.avatar || (props.profile.first_name ? props.profile.first_name[0].toUpperCase() : '')}
    </Avatar>
  );
}
