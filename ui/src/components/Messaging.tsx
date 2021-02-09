import { SendOutlined } from '@ant-design/icons';
import { Layout, Input } from 'antd';
import React, { useState } from 'react';
import styled from 'styled-components';
import { ReadProfileDTO } from '../api/definitions';
import { format } from 'date-and-time';

const StyledLayout = styled(Layout)`
  width: 100%;
  background-color: rgb(10 10 10);
`;

const StyledMessages = styled.div`
  display: flex;
  flex-direction: column;
  height: calc(100% - 32px);
  padding: 1em;
  overflow-y: scroll;
  &::-webkit-scrollbar {
    width: 3px;
  }

  /* Track */
  &::-webkit-scrollbar-track {
    background: #080808;
  }

  /* Handle */
  &::-webkit-scrollbar-thumb {
    background: #2b2b2b;
  }

  /* Handle on hover */
  &::-webkit-scrollbar-thumb:hover {
    background: #555;
  }
`;

const { Search } = Input;

const StyledTextArea = styled(Search)`
  position: absolute;
  bottom: 0;
  height: 32px;
  & input {
    background-color: rgb(5 5 5);
    border: 1px solid rgb(5 5 5);
    color: #fff;
  }
  & input::placeholder {
    color: #3d3d3d;
  }
`;

const StyledMessage = styled.div<{ self: boolean }>`
  background: ${(props) => (props.self ? '#2d2d2d' : '#1890ff')};
  border-radius: 0.3em;
  padding: 0.1em 1em;
  max-width: 80%;
  width: fit-content;
  margin-bottom: 1em;
  display: block;
  clear: both;
  ${(props) => (props.self ? 'align-self: flex-end; text-align: right;' : 'align-self: flex-start;')}
  color: #fff;
  & span {
    color: #818181;
    display: block;
    text-align: right;
    font-size: 0.7em;
  }
`;

interface Message {
  profile?: ReadProfileDTO;
  date: Date;
  text: string;
}

interface MessagingProps {
  height: number;
}

export default function Messaging(props: MessagingProps): JSX.Element {
  const [messages, setMessages] = useState<Message[]>([]);
  const [text, setText] = useState('');
  const sendMessage = () => {
    if (text) {
      setMessages(messages.concat({ text, date: new Date() }));
      setText('');
    }
  };
  return (
    <StyledLayout style={{ height: props.height }}>
      <StyledMessages>
        {messages.map((v, i) => (
          <StyledMessage key={i} self={!v.profile}>
            {v.text}
            <span>{format(v.date, 'h:mm A')}</span>
          </StyledMessage>
        ))}
      </StyledMessages>
      <StyledTextArea
        placeholder="Send a Message"
        value={text}
        onChange={(e) => setText(e.currentTarget.value)}
        onSearch={sendMessage}
        enterButton={<SendOutlined />}
      ></StyledTextArea>
    </StyledLayout>
  );
}
