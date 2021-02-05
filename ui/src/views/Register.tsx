
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
  MailOutlined
} from "@ant-design/icons";

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const StyledLayout = styled(Layout)`
  padding: 2em 0;
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

  return (
    <StyledLayout>
      <Row justify="center">
        <Col span={8}>
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
              <Button type="primary" htmlType="submit">Register</Button>
              <Button type="link" href="/login">Log in</Button>
            </Form.Item>
          </Form>
        </Col>
      </Row>
    </StyledLayout>
  );
}



export default Register;