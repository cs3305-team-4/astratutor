import { Layout } from 'antd';
import React from 'react';
import styled from 'styled-components';

const StyledLayout = styled(Layout)`
  width: 100%;
  background-color: rgb(10 10 10);
`;

interface MessagingProps {
  height: number;
}

export default function Messaging(props: MessagingProps): JSX.Element {
  return <StyledLayout style={{ height: props.height }}></StyledLayout>;
}
