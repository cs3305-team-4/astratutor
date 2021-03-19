import { SendOutlined, SettingFilled } from '@ant-design/icons';
import { Layout, Input } from 'antd';
import React, { useEffect, useRef, useState } from 'react';
import styled from 'styled-components';
import { ProfileResponseDTO } from '../api/definitions';
import { UserAvatar } from './UserAvatar';

const StyledLayout = styled(Layout)`
  width: 100%;
  display: flex;
  flex: 1;
  background-color: rgb(10 10 10);
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

const StyledMessages = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: 15px 1em 32px;
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
  word-wrap: break-word;
  border-radius: 0.3em;
  padding: 0.1em 1em;
  max-width: ${(props) => (props.self ? '80%' : '60%')};
  width: fit-content;
  margin-bottom: 1em;
  display: block;
  clear: both;
  ${(props) => (props.self ? 'align-self: flex-end; text-align: right;' : 'align-self: flex-start;')}
  color: #fff;
  & span {
    color: ${(props) => (props.self ? '#818181' : '#e2e2e2')};
    display: block;
    text-align: right;
    font-size: 0.7em;
  }
`;

export interface Message {
  profile?: ProfileResponseDTO;
  date: Date;
  text: string;
}

interface MessagingProps {
  height: number;
  messages: Message[];
  setMessages: (m: Message[]) => void;
}

export default function Messaging(props: MessagingProps): JSX.Element {
  const [text, setText] = useState('');
  const [last, setLast] = useState<Message>();

  const sendMessage = () => {
    if (text) {
      props.setMessages(props.messages.concat({ text, date: new Date() }));
      setText('');
    }
  };

  useEffect(() => {
    const el = document.getElementById('messages');
    if (el) {
      el.scrollTop = el.scrollHeight;
    }
  }, [last]);
  return (
    <StyledLayout id="messages" style={{ height: `calc(100vh - ${props.height}px)` }}>
      <div style={{ flexGrow: 2 }}></div>
      <StyledMessages onChange={(e) => console.log('resize')}>
        {props.messages.map((v, i) => {
          if (i === props.messages.length - 1 && last?.date !== v.date) {
            setLast(v);
          }
          if (v.profile) {
            return (
              <div key={i} style={{ display: 'flex' }}>
                {v.profile && <UserAvatar profile={v.profile} props={{ style: { float: 'left' } }}></UserAvatar>}
                <StyledMessage self={!v.profile} style={{ float: 'left', marginLeft: '1em' }}>
                  {v.text}
                  <span>
                    {new Intl.DateTimeFormat('en-IE', {
                      hour: 'numeric',
                      minute: 'numeric',
                    }).format(v.date)}
                  </span>
                </StyledMessage>
              </div>
            );
          }
          return (
            <StyledMessage key={i} self={!v.profile}>
              {v.text}
              <span style={{ wordWrap: 'break-word' }}>
                {new Intl.DateTimeFormat('en-IE', {
                  hour: 'numeric',
                  minute: 'numeric',
                }).format(v.date)}
              </span>
            </StyledMessage>
          );
        })}
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
