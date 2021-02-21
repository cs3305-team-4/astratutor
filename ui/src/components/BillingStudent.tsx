import React from 'react';
import styled from 'styled-components';

import { useAsync } from 'react-async-hook';

import {
  Typography,
  Layout,
  Row,
  Col,
  Avatar,
  PageHeader,
  Input,
  Button,
  Statistic,
  Form,
  Upload,
  Modal,
  Alert,
  Table,
  Skeleton,
} from 'antd';

import { BankOutlined } from '@ant-design/icons';

import { AccountType, ProfileRequestDTO } from '../api/definitions';
import { APIContext } from '../api/api';
import DefaultAvatar from '../assets/default_avatar.png';

import { Redirect, useHistory } from 'react-router-dom';

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;
const { TextArea } = Input;

export function BillingStudent(): React.ReactElement {
  const [error, setError] = React.useState<string>('');

  const [billingRequirementsMet, setBillingRequirementsMet] = React.useState<boolean>(true);
  const [ready, setReady] = React.useState<boolean>(false);
  const api = React.useContext(APIContext);
  const history = useHistory();

  useAsync(async () => {
    setReady(true);
    // const onboard = await api.services.readTutorBillingOnboardStatus(api.account.id);
    // setBillingOnboard(onboard);
    // if (!onboard) {
    //   setBillingOnboardUrl(await api.services.readTutorBillingOnboardUrl(api.account.id));
    //   setReady(true);
    // } else {
    //   const reqsMet = await api.services.readTutorBillingRequirementsMetStatus(api.account.id);
    //   setBillingRequirementsMet(reqsMet);
    //   setBillingPanelUrl(await api.services.readTutorBillingPanelUrl(api.account.id));
    //   setReady(true);
    // }
  }, []);

  const redirectBillingOnboard = () => {
    window.location.href = billingOnboardUrl;
  };

  const redirectBillingPanelAccount = () => {
    window.open(billingPanelUrl + '#/account');
  };

  const redirectBillingPanel = () => {
    window.open(billingPanelUrl + '#/account');
  };

  if (!ready) {
    return (
      <>
        <Skeleton />
      </>
    );
  }

  return (
    <Typography>
      <Content>
        <Row gutter={16} style={{ margin: '1rem' }}>
          <Title level={3}>
            <BankOutlined /> Billing
          </Title>
          <Col md={24} sm={24} xs={24}>
            {!false && (
              <>
                <Alert
                  message="No Debit/Credit Card Available"
                  description="You do not have a debit or credit card linked to your account, you will need to link a card when requesting a lesson"
                  type="warning"
                  showIcon
                  style={{ margin: '1rem 0' }}
                />
              </>
            )}
          </Col>
          <Col md={24} sm={24} xs={24}>
            {!false && (
              <>
                <Button style={{ margin: '0.5em' }} onClick={redirectBillingPanelAccount}>
                  Modify Card Details
                </Button>
              </>
            )}
          </Col>
        </Row>
        <Row gutter={16} style={{ margin: '1rem' }}>
          <Title level={3}>Invoices</Title>
          <Col md={24} sm={24} xs={24}>
            <Table
              locale={{
                emptyText: 'No invoices available',
              }}
              columns={[
                { title: 'Lesson', key: 'degree', dataIndex: 'degree' },
                { title: 'Date', key: 'field', dataIndex: 'field' },
                { title: 'Amount', key: 'school', dataIndex: 'school' },
                { title: 'Available for Payout', key: 'verified', dataIndex: 'verified' },
                { title: 'Paid Out', key: 'verified', dataIndex: 'verified' },
                { title: '', key: 'delete', dataIndex: 'delete' },
              ]}
              size="small"
              style={{ width: '100%' }}
              pagination={false}
              // dataSource={profile.qualifications.map((quali: QualificationResponseDTO) => {
              //   return {
              //     degree: quali.degree,
              //     field: quali.field,
              //     school: quali.school,
              //     verified: quali.verified ? '\u2713' : '\u2717',
              //     delete: editQualis ? (
              //       <Button onClick={() => deleteQuali(quali.id)} style={{ margin: '0 0.5rem' }} size="small">
              //         <DeleteOutlined />
              //         Remove
              //       </Button>
              //     ) : (
              //       <></>
              //     ),
              // //   };
              // })}
            ></Table>
          </Col>
        </Row>
      </Content>
    </Typography>
  );
}
