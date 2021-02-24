import { ProfileResponseDTO } from '../api/definitions';
import { Avatar, AvatarProps, Tooltip } from 'antd';
import React from 'react';

export function UserAvatar(props: { profile: ProfileResponseDTO; props?: AvatarProps }): JSX.Element {
  return (
    <Tooltip title={props.profile.first_name + ' ' + props.profile.last_name} placement="top">
      {props.profile.avatar ? (
        <Avatar
          {...props.props}
          style={{ backgroundColor: props.profile.color, cursor: 'pointer', ...props.props?.style }}
          icon={<img src={props.profile.avatar} alt="profile" />}
        ></Avatar>
      ) : (
        <Avatar
          {...props.props}
          style={{ backgroundColor: props.profile.color, cursor: 'pointer', ...props.props?.style }}
        >
          {props.profile.avatar || (props.profile.first_name ? props.profile.first_name[0].toUpperCase() : '')}
        </Avatar>
      )}
    </Tooltip>
  );
}
