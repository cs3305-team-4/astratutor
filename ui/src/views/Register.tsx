
import React, { useReducer, useState } from 'react';
import styled from 'styled-components';

import {
  Layout,
  Typography,
  Row,
  Col,
  Form, 
  Input, 
  Button, 
  Checkbox,
  Radio,
  Space
} from "antd";
import {
  LockOutlined,
  MailOutlined,
  UserAddOutlined,
} from "@ant-design/icons";

import { AuthContext } from '../api/auth'
import DeskImg from "../assets/stock/desk-medium.jpg"

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const StyledLayout = styled(Layout)`
  height: calc(100vh - 72px);
  background-image: url(${DeskImg});
  background-size: cover;
`;

type RegisterState = {
  accountType: 'student' | 'tutor';
  under16: boolean;
}

type RegisterStateAction =
  | { type: "account-student" }
  | { type: "account-tutor" }
  | { type: "toggle-under-16" };

const Register: React.FunctionComponent = () => {
  const initialState: RegisterState = {
    accountType: 'student',
    under16: false
  };
  
  const reducer = (state: RegisterState, action: RegisterStateAction): RegisterState => {
    switch(action.type) {
      case "account-student":
        return { accountType: 'student', under16: state.under16 };
      case "account-tutor":
        state.accountType = "tutor";
        return { accountType: 'tutor', under16: state.under16 };
      case "toggle-under-16":
        return { accountType: state.accountType, under16: !state.under16 };
    }
  }

  const [state, dispatch] = useReducer(reducer, initialState);

  const onSubmit = (values: any) => {
    console.log("Success:", values);
  };

  const FormItems = () => {
    // All default items in a register form
    const items = [
      <Form.Item
        name="email"
        rules={[{required: true, message: 'Please input your email!'}]}
      >
        <Input
          prefix={<MailOutlined />}
          placeholder="Email" 
        />
      </Form.Item>,
      <Form.Item
        name="password"
        rules={[{required: true, message: 'Please enter your password!'}]}
      >
        <Input.Password
          prefix={<LockOutlined />}
          placeholder="Password" 
        />
      </Form.Item>,
      <Form.Item
        rules={[{required: true, message: 'Please confirm your password!'}]}
      >
        <Input.Password
          prefix={<LockOutlined />}
          placeholder="Confirm Password"
        />
      </Form.Item>
    ];

    // Items specific to student registration
    if (state.accountType == 'student') {
      items.push(
        <Form.Item
          name="under16"
          valuePropName="checked"
          initialValue={state.under16}
        >
          <Checkbox 
            onChange={() => dispatch({ type: "toggle-under-16" })}
            value={state.under16}
          >
            I am under the age of 16
          </Checkbox>
        </Form.Item>
      );

      if (state.under16) {
        items.push(
          <Form.Item
            name="parentsEmail"
          >
            <Input
              prefix={<MailOutlined />}
              placeholder="Parents Email"
            />
          </Form.Item>
        );
      }
    }

    // Accept Terms
    items.push(
      <Form.Item
        name="terms"
        valuePropName="checked"
        initialValue={false}
      >
        <Checkbox>I accept the <a href="/terms">Terms and Conditions</a></Checkbox>
      </Form.Item>
    )

    return items; 
  }


  const authContext = React.useContext(AuthContext)
  if (authContext.isLoggedIn()) {
    return (
      <StyledLayout>
        <Content>
          <Row
            style={{height: "calc(100vh - 72px)"}}
            align="middle"
            justify="center"
          >
            <Col md={9} sm={6} xs={0}/>
            <Col md={6} sm={10} xs={24} style={{padding: "1rem", backgroundColor: "rgba(255,255,255,0.8)"}}>
              You are already logged in. Click <Link to="/">here</Link> to continue
            </Col>
            <Col md={9} sm={6} xs={0}/>
          </Row>
        </Content>
      </StyledLayout>
    )
  }

  return (
    <StyledLayout>
      <Row
          style={{height: "calc(100vh - 72px)"}}
          align="middle"
          justify="center"
      >
        <Col md={10} sm={6} xs={0}/>
        <Col md={4} sm={10} xs={24} style={{padding: "1rem", backgroundColor: "rgba(255,255,255,0.8)"}}>
          <UserAddOutlined style={{ display: "block", margin: "0 auto", fontSize: "4rem", padding: "2rem", color: "rgb(200,200,200)"}} />
          <Form
            onFinish={onSubmit}
          >
            <Form.Item>
              <Row justify="center">
                <Radio.Group value={state.accountType}>
                  <Radio.Button 
                    onClick={() => dispatch({ type: "account-student" })}
                    value="student"
                  >
                    Student
                  </Radio.Button>
                  <Radio.Button 
                    onClick={() => dispatch({ type: "account-tutor" })}
                    value="tutor"
                  >
                    Tutor
                  </Radio.Button>
                </Radio.Group>
              </Row>
            </Form.Item>

            {FormItems()}
            
            <Form.Item>
              <Button style={{width: "100%"}} type="primary" htmlType="submit">
                Register
              </Button> 
            </Form.Item>
          </Form>
        </Col>
        <Col md={10} sm={6} xs={0}/>
      </Row>
    </StyledLayout>
  );
}



export default Register;