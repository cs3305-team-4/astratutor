import React from 'react';
import styled from 'styled-components';

import { useAsync } from 'react-async-hook';

import {
  Alert,
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
  Select,
  Modal,
  Tabs,
  Tag,
  Skeleton,
  InputNumber,
  AvatarProps,
  Table,
  Progress,
} from 'antd';

import { UploadRequestOption } from 'rc-upload/lib/interface';
import { useHistory } from 'react-router';

import ImgCrop from 'antd-img-crop';

import {
  EditOutlined,
  UserAddOutlined,
  PlusOutlined,
  UserOutlined,
  DeleteOutlined,
  CheckOutlined,
  StopOutlined,
} from '@ant-design/icons';

import {
  AccountType,
  ProfileRequestDTO,
  SubjectTaughtRequestDTO,
  SubjectTaughtPriceUpdateRequestDTO,
  SubjectTaughtDescriptionUpdateRequestDTO,
  WorkExperienceRequestDTO,
  WorkExperienceResponseDTO,
  QualificationRequestDTO,
  QualificationResponseDTO,
  ProfileResponseDTO,
  SubjectTaughtDTO,
  SubjectDTO,
} from '../api/definitions';

import { RequestLessonModal } from './RequestLessonModal';

import { APIContext } from '../api/api';
import { Availability } from './Availability';
import { UserAvatar } from './UserAvatar';

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const { TextArea } = Input;

export interface ProfileProps {
  uuid: string;
  type: AccountType;
}

export function Profile(props: ProfileProps): React.ReactElement {
  const api = React.useContext(APIContext);
  const history = useHistory();

  const [profile, setProfile] = React.useState<ProfileResponseDTO | undefined>(undefined);
  const [activeTab, setActiveTab] = React.useState<string>('outline');
  const [tutorSubjects, setTutorSubjects] = React.useState<SubjectTaughtDTO[] | undefined>(undefined);
  const [subjects, setSubjects] = React.useState<SubjectDTO[] | undefined>(undefined);
  const [TutorSubjectID, setTutorSubjectID] = React.useState<string>();

  const isSelf: boolean = api.isLoggedIn() && props.uuid === api.account.id;

  const [requestLessonVisible, setRequestLessonVisible] = React.useState<boolean>(false);

  const [editSubtitle, setEditSubtitle] = React.useState<boolean>(false);
  const [newSubtitle, setNewSubtitle] = React.useState<string>('');
  const [addSubVisible, setAddSubVisible] = React.useState<boolean>(false);

  const [editDesc, setEditDesc] = React.useState<boolean>(false);
  const [newDesc, setNewDesc] = React.useState<string>();

  const [editSubs, setEditSubs] = React.useState<boolean>(false);
  const [editSubDescVisible, setEditSubDescVisible] = React.useState<boolean>(false);
  const [editSubPriceVisible, setEditSubPriceVisible] = React.useState<boolean>(false);

  const [editQualis, setEditQualis] = React.useState<boolean>(false);
  const [addQualiVisible, setAddQualiVisible] = React.useState<boolean>(false);

  const [editAvail, setEditAvail] = React.useState<boolean>(false);

  const [editWork, setEditWork] = React.useState<boolean>(false);
  const [addWorkVisible, setAddWorkVisible] = React.useState<boolean>(false);

  const reloadProfile = async () => {
    try {
      setProfile(await api.services.readProfileByAccountID(props.uuid, props.type));
      setTutorSubjects(await api.services.readTutorSubjectsByAccountId(props.uuid));
      setSubjects(await api.services.readSubjects());
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not load profile: ${e}`,
      });
    }
  };

  useAsync(async () => {
    await reloadProfile();
  }, []);

  const commitHours = async (hours: boolean[]) => {
    try {
      await api.services.updateAvailabilityOnProfileID(profile.account_id, AccountType.Tutor, hours);
      await reloadProfile();
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not set availability: ${e}`,
      });
    }
  };

  const commitQuali = async (quali: QualificationRequestDTO) => {
    try {
      await api.services.createQualificationOnProfileID(props.uuid, props.type, quali);
      await reloadProfile();
      setAddQualiVisible(false);
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not create qualification: ${e}`,
      });
    }
  };

  const deleteQuali = async (id: string) => {
    try {
      await api.services.deleteQualificationOnProfileID(props.uuid, props.type, id);
      await reloadProfile();
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not delete qualification: ${e}`,
      });
    }
  };

  const subPricEdit = async (id: string) => {
    setTutorSubjectID(id);
    setEditSubPriceVisible(!editSubDescVisible);
  };

  const subDescEdit = async (id: string) => {
    setTutorSubjectID(id);
    setEditSubDescVisible(!editSubDescVisible);
  };

  const commitWork = async (work: WorkExperienceRequestDTO) => {
    try {
      await api.services.createWorkExperienceOnProfileID(props.uuid, props.type, work);
      await reloadProfile();
      setAddQualiVisible(false);
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not create work experience: ${e}`,
      });
    }
  };

  const commitSub = async (subjectTaught: SubjectTaughtRequestDTO) => {
    console.log(subjectTaught);
    try {
      await api.services.createSubjectTaughtOnProfileID(props.uuid, props.type, subjectTaught);
      await reloadProfile();
      setAddSubVisible(false);
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not teach Subject: ${e}`,
      });
    }
  };

  const deleteWork = async (id: string) => {
    try {
      await api.services.deleteWorkExperienceOnProfileID(props.uuid, props.type, id);
      await reloadProfile();
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not delete work experience: ${e}`,
      });
    }
  };

  const commitDesc = async (desc: string) => {
    try {
      await api.services.updateDescriptionOnProfileID(props.uuid, props.type, desc);
      await reloadProfile();
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not set description: ${e}`,
      });
    }
  };

  const commitSubtitle = async (Subtitle: string) => {
    try {
      await api.services.updateSubtitleOnProfileID(props.uuid, props.type, Subtitle);
      await reloadProfile();
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not set subtitle: ${e}`,
      });
    }
  };

  const commitSubPrice = async (price: SubjectTaughtPriceUpdateRequestDTO) => {
    try {
      await api.services.updateSubjectPriceOnProfileID(props.uuid, TutorSubjectID, props.type, price);
      await reloadProfile();
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not set description: ${e}`,
      });
    }
  };

  const commitSubDescription = async (desc: SubjectTaughtDescriptionUpdateRequestDTO) => {
    console.log(desc);
    try {
      await api.services.updateSubjectDescriptionOnProfileID(props.uuid, TutorSubjectID, props.type, desc);
      await reloadProfile();
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not set description: ${e}`,
      });
    }
  };

  const commitAvatar = async (opt: UploadRequestOption) => {
    const reader = new FileReader();
    reader.readAsDataURL(opt.file);
    reader.onload = async (e) => {
      try {
        const img = new Image();

        img.onload = async (e) => {
          const canvas = document.createElement('canvas');
          canvas.width = 256;
          canvas.height = 256;
          const ctx = canvas.getContext('2d');
          ctx.drawImage(img, 0, 0, canvas.width, canvas.height);

          await api.services.updateAvatarOnProfileID(props.uuid, props.type, canvas.toDataURL('image/jpeg', 0.8));
          await reloadProfile();
        };
        img.src = reader.result as string;
      } catch (e) {
        Modal.error({
          title: 'Error',
          content: `Could not set avatar: ${e}`,
        });
      }
    };

    reader.onerror = (error: ProgressEvent<FileReader>) => {
      throw error;
    };
  };

  if (profile === undefined) {
    return (
      <>
        <Skeleton />
      </>
    );
  }

  return (
    <Typography>
      {isSelf && (
        <Alert
          style={{ margin: '1rem' }}
          message="Your Profile"
          description="Your profile will present like this to others (without the option of editing elements)"
          type="info"
          showIcon
        />
      )}
      <PageHeader
        title={
          <>
            <div style={{ padding: '0.5rem 0' }}>
              <ImgCrop
                rotate
                shape="round"
                cropperProps={{
                  cropSize: {
                    width: 128,
                    height: 128,
                  },
                }}
              >
                <Upload
                  defaultFileList={[]}
                  style={{ borderRadius: '16px', width: '64px', display: 'inline' }}
                  // onChange={onChange}
                  // onPreview={onPreview}
                  fileList={[]}
                  customRequest={commitAvatar}
                  disabled={!isSelf}
                >
                  <Button
                    style={{
                      border: 'none',
                      padding: '0',
                      margin: '0',
                      zIndex: 10,
                      width: '96px',
                      height: '96px',
                      position: 'relative',
                    }}
                    size="small"
                  >
                    <UserAvatar profile={profile} props={{ size: 96, style: { fontSize: 40 } }} />
                  </Button>
                </Upload>
              </ImgCrop>
            </div>
            {`${profile.first_name} ${profile.last_name}`}
            <Title level={5} style={{ fontWeight: 300, margin: '0 0 0.5rem 0' }}>
              <Button
                hidden={!isSelf}
                onClick={async () => {
                  if (editSubtitle === false) {
                    setNewSubtitle(profile.subtitle);
                  } else {
                    await commitSubtitle(newSubtitle);
                  }

                  setEditSubtitle(editSubtitle ? false : true);
                }}
                size="small"
                style={{ margin: '0 0.5rem' }}
                type={editSubtitle ? 'primary' : 'default'}
              >
                <EditOutlined />
                {!editSubtitle ? 'Edit' : 'Finish'}
              </Button>
              {!editSubtitle ? (
                <Paragraph style={{ whiteSpace: 'pre-wrap' }}>{profile.subtitle}</Paragraph>
              ) : (
                <input
                  maxLength={300}
                  onChange={(ev) => {
                    setNewSubtitle(ev.target.value);
                  }}
                  style={{ minHeight: '240px', margin: '0.5rem 0' }}
                  value={newSubtitle}
                />
              )}
            </Title>
            {props.type === AccountType.Tutor && (
              <Content style={{ display: 'flex', flexWrap: 'wrap' }}>
                {tutorSubjects?.map((subject, index) => (
                  <Tag style={{ margin: '0.25rem' }} key={index}>
                    {subject.name}
                  </Tag>
                ))}
              </Content>
            )}
          </>
        }
        className="site-page-header-responsive"
        extra={[
          <Row key="stats" style={{ marginTop: '2.5rem' }} align="top" justify="start" gutter={8}>
            <Col>
              <Statistic key="users" title="Location" value={`${profile.city}, ${profile.country}`} />
            </Col>
            <Col>
              <Statistic key="users" title="User Since" value={'March 28th 2019'} />
            </Col>
            {props.type === AccountType.Tutor && (
              <Col>
                <Statistic key="users" title="Lessons Given" value={24} />
              </Col>
            )}
            {props.type === AccountType.Tutor && (
              <Col>
                <Statistic key="users" title="Average Review" value={'4.5/5'} />
              </Col>
            )}
          </Row>,
          <Row key="buttons" gutter={16} align="top" justify="end" style={{ margin: '0.5rem 0' }}>
            <Button
              type="primary"
              key="request"
              style={{ margin: '0.2rem' }}
              onClick={() => {
                if (!api.isLoggedIn()) {
                  history.push('/login');
                } else {
                  setRequestLessonVisible(true);
                }
              }}
              disabled={isSelf}
            >
              Request Lesson
            </Button>
          </Row>,
        ]}
        footer={
          props.type === AccountType.Tutor && (
            <Tabs defaultActiveKey={activeTab} onChange={(key: string) => setActiveTab(key)}>
              <Tabs.TabPane tab="Outline" key="outline" />
              <Tabs.TabPane tab="Reviews" key="reviews" />
            </Tabs>
          )
        }
        style={{ padding: '0.1rem 1rem' }}
      ></PageHeader>
      {props.type === AccountType.Tutor && (
        <Content style={{ padding: '1rem' }}>
          {activeTab === 'outline' && (
            <>
              <Row gutter={16}>
                <Col md={12} sm={24} xs={24} style={{ margin: '0.5rem 0' }}>
                  <Title level={5}>
                    Description
                    <Button
                      hidden={!isSelf}
                      onClick={async () => {
                        if (editDesc === false) {
                          setNewDesc(profile.description);
                        } else {
                          await commitDesc(newDesc);
                        }

                        setEditDesc(editDesc ? false : true);
                      }}
                      size="small"
                      style={{ margin: '0 0.5rem' }}
                      type={editDesc ? 'primary' : 'default'}
                    >
                      <EditOutlined />
                      {!editDesc ? 'Edit' : 'Finish'}
                    </Button>
                  </Title>
                  {!editDesc ? (
                    <Paragraph style={{ whiteSpace: 'pre-wrap' }}>{profile.description}</Paragraph>
                  ) : (
                    <TextArea
                      maxLength={1000}
                      onChange={(ev) => {
                        setNewDesc(ev.target.value);
                      }}
                      style={{ minHeight: '240px', margin: '0.5rem 0' }}
                      value={newDesc}
                      size="large"
                    />
                  )}
                  <Title level={5}>
                    Subjects
                    <Button
                      hidden={!isSelf}
                      onClick={async () => {
                        setEditSubs(editSubs ? false : true);
                      }}
                      size="small"
                      style={{ margin: '0 0.5rem' }}
                      type={editSubs ? 'primary' : 'default'}
                    >
                      <EditOutlined />
                      {!editSubs ? 'Edit' : 'Finish'}
                    </Button>
                    {editSubs && (
                      <Button
                        size="small"
                        style={{ position: 'relative', right: 0, margin: '0 0.5rem' }}
                        onClick={() => setAddSubVisible(!addSubVisible)}
                      >
                        <PlusOutlined />
                        Teach A New Subject
                      </Button>
                    )}
                  </Title>
                  <Table
                    locale={{
                      emptyText: 'No Subjects listed',
                    }}
                    columns={[
                      { title: 'Subject', key: 'subject', dataIndex: 'subject' },
                      { title: 'Price', key: 'price', dataIndex: 'price' },
                      { title: 'Description', key: 'description', dataIndex: 'description' },
                      { title: '', key: 'editPrice', dataIndex: 'editPrice' },
                      { title: '', key: 'editDesc', dataIndex: 'editDesc' },
                    ]}
                    size="small"
                    style={{ width: '100%' }}
                    pagination={false}
                    dataSource={tutorSubjects?.map((subject, index) => {
                      return {
                        price: subject.price,
                        subject: subject.name,
                        description: subject.description,
                        editDesc: editSubs ? (
                          <Button onClick={() => subDescEdit(subject.id)} style={{ margin: '0 0.5rem' }} size="small">
                            <EditOutlined />
                            Edit Description
                          </Button>
                        ) : (
                          <></>
                        ),
                        editPrice: editSubs ? (
                          <Button onClick={() => subPricEdit(subject.id)} style={{ margin: '0 0.5rem' }} size="small">
                            <EditOutlined />
                            Edit Price
                          </Button>
                        ) : (
                          <></>
                        ),
                      };
                    })}
                  ></Table>
                  <Row style={{ margin: '0.5rem 0' }}>
                    <Title level={5}>
                      Qualifications
                      <Button
                        hidden={!isSelf}
                        onClick={async () => {
                          setEditQualis(editQualis ? false : true);
                        }}
                        size="small"
                        style={{ margin: '0 0.5rem' }}
                        type={editQualis ? 'primary' : 'default'}
                      >
                        <EditOutlined />
                        {!editQualis ? 'Edit' : 'Finish'}
                      </Button>
                      {editQualis && (
                        <Button
                          size="small"
                          style={{ position: 'relative', right: 0, margin: '0 0.5rem' }}
                          onClick={() => setAddQualiVisible(!addQualiVisible)}
                        >
                          <PlusOutlined />
                          Add
                        </Button>
                      )}
                    </Title>
                    <Table
                      locale={{
                        emptyText: 'No qualifications listed',
                      }}
                      columns={[
                        { title: 'Degree', key: 'degree', dataIndex: 'degree' },
                        { title: 'Field', key: 'field', dataIndex: 'field' },
                        { title: 'Awarding Institution', key: 'school', dataIndex: 'school' },
                        { title: 'Verified', key: 'verified', dataIndex: 'verified' },
                        { title: '', key: 'delete', dataIndex: 'delete' },
                      ]}
                      size="small"
                      style={{ width: '100%' }}
                      pagination={false}
                      dataSource={profile.qualifications.map((quali: QualificationResponseDTO) => {
                        return {
                          degree: quali.degree,
                          field: quali.field,
                          school: quali.school,
                          verified: quali.verified ? '\u2713' : '\u2717',
                          delete: editQualis ? (
                            <Button onClick={() => deleteQuali(quali.id)} style={{ margin: '0 0.5rem' }} size="small">
                              <DeleteOutlined />
                              Remove
                            </Button>
                          ) : (
                            <></>
                          ),
                        };
                      })}
                    ></Table>
                  </Row>
                  <Row style={{ margin: '0.5rem 0' }}>
                    <Title level={5}>
                      Work Experience
                      <Button
                        hidden={!isSelf}
                        onClick={async () => {
                          setEditWork(editWork ? false : true);
                        }}
                        size="small"
                        style={{ margin: '0 0.5rem' }}
                        type={editWork ? 'primary' : 'default'}
                      >
                        <EditOutlined />
                        {!editWork ? 'Edit' : 'Finish'}
                      </Button>
                      {editWork && (
                        <Button
                          size="small"
                          style={{ position: 'relative', right: 0, margin: '0 0.5rem' }}
                          onClick={() => setAddWorkVisible(!addWorkVisible)}
                        >
                          <PlusOutlined />
                          Add
                        </Button>
                      )}
                    </Title>
                    <Table
                      locale={{
                        emptyText: 'No work experience listed',
                      }}
                      columns={[
                        { title: 'Role', key: 'role', dataIndex: 'role' },
                        { title: 'Years Exp.', key: 'years_exp', dataIndex: 'years_exp' },
                        { title: 'Description', key: 'description', dataIndex: 'description' },
                        { title: 'Verified', key: 'verified', dataIndex: 'verified' },
                        { title: '', key: 'delete', dataIndex: 'delete' },
                      ]}
                      size="small"
                      style={{ width: '100%' }}
                      pagination={false}
                      dataSource={profile.work_experience.map((exp: WorkExperienceResponseDTO) => {
                        return {
                          role: exp.role,
                          years_exp: exp.years_exp,
                          description: exp.description,
                          verified: exp.verified ? '\u2713' : '\u2717',
                          delete: editWork ? (
                            <Button onClick={() => deleteWork(exp.id)} style={{ margin: '0 0.5rem' }} size="small">
                              <DeleteOutlined />
                              Remove
                            </Button>
                          ) : (
                            <></>
                          ),
                        };
                      })}
                    ></Table>
                  </Row>
                </Col>
                <Col md={12} sm={24} xs={24} style={{ margin: '0.5rem 0' }}>
                  <Title level={5}>
                    Availability
                    <Button
                      hidden={!isSelf}
                      size="small"
                      style={{ margin: '0 0.5rem' }}
                      onClick={() => setEditAvail(!editAvail)}
                      type={editAvail ? 'primary' : 'default'}
                    >
                      <EditOutlined />
                      {!editAvail ? 'Edit' : 'Finish'}
                    </Button>
                  </Title>
                  <Availability
                    hideUnavailable={true}
                    availability={profile.availability}
                    onChange={commitHours}
                    editable={editAvail}
                  ></Availability>
                </Col>
              </Row>
              <Modal
                title="Add Qualification"
                visible={addQualiVisible}
                onCancel={() => setAddQualiVisible(false)}
                footer={[
                  <Button form="add-quali" key="submit" style={{ width: '100%' }} type="primary" htmlType="submit">
                    Add
                  </Button>,
                ]}
              >
                <Form
                  onFinish={commitQuali}
                  initialValues={{ degree: 'Bachelors' }}
                  layout="vertical"
                  name="add-quali"
                  preserve={false}
                >
                  <Form.Item name="degree" rules={[{ required: true, message: 'Please select a degree type!' }]}>
                    <Select size="large" style={{ width: '100%' }}>
                      {['Associates', 'Bachelors', 'Masters', 'Doctorate'].map((value, index) => (
                        <Select.Option key={value} value={value}>
                          {value} degree
                        </Select.Option>
                      ))}
                    </Select>
                  </Form.Item>
                  <Form.Item name="field" rules={[{ required: true, message: 'Please name your field!' }]}>
                    <Input
                      size="large"
                      placeholder="Field (i.e Biology, Computer Science, Arts)"
                      style={{ width: '100%' }}
                    ></Input>
                  </Form.Item>
                  <Form.Item
                    name="school"
                    rules={[{ required: true, message: 'Please name the awarding institution!' }]}
                  >
                    <Input size="large" placeholder="Awarding Institution" style={{ width: '100%' }}></Input>
                  </Form.Item>
                </Form>
              </Modal>
              <Modal
                title="Add Work Experience"
                visible={addWorkVisible}
                onCancel={() => setAddWorkVisible(false)}
                footer={[
                  <Button form="add-work" key="submit" style={{ width: '100%' }} type="primary" htmlType="submit">
                    Add
                  </Button>,
                ]}
              >
                <Form onFinish={commitWork} layout="vertical" name="add-work" preserve={false}>
                  <Form.Item name="role" rules={[{ required: true, message: 'Please enter a role!' }]}>
                    <Input size="large" placeholder="Role" style={{ width: '100%' }}></Input>
                  </Form.Item>
                  <Form.Item
                    name="years_exp"
                    rules={[{ required: true, message: 'Please provide the number of years experience!' }]}
                  >
                    <InputNumber
                      placeholder="Number of years experience"
                      size="large"
                      style={{ width: '100%' }}
                      min={1}
                      max={50}
                    />
                  </Form.Item>
                  <Form.Item name="description" rules={[{ required: true, message: 'Please enter a description!' }]}>
                    <TextArea
                      maxLength={240}
                      placeholder="Short description of the role"
                      style={{ minHeight: '240px', margin: '0.5rem 0' }}
                      size="large"
                    />
                  </Form.Item>
                </Form>
              </Modal>

              <Modal
                title="Teach a Subject"
                visible={addSubVisible}
                onCancel={() => setAddSubVisible(false)}
                footer={[
                  <Button form="add-subject" key="submit" style={{ width: '100%' }} type="primary" htmlType="submit">
                    Add
                  </Button>,
                ]}
              >
                <Form onFinish={commitSub} layout="vertical" name="add-subject" preserve={false}>
                  <Form.Item name="subject_id" rules={[{ required: true, message: 'Please select Subject!' }]}>
                    <Select size="large" showSearch style={{ width: '100%' }}>
                      {subjects?.map((subject, index) => (
                        <Select.Option key={index} value={subject.id}>
                          {subject.name}
                        </Select.Option>
                      ))}
                    </Select>
                  </Form.Item>
                  <Form.Item
                    name="price"
                    rules={[{ required: true, message: 'Please provide how much you wish your subject to cost.' }]}
                  >
                    <InputNumber
                      placeholder="Desired Cost of a Lesson"
                      size="large"
                      style={{ width: '100%' }}
                      min={1}
                    />
                  </Form.Item>
                  <Form.Item name="description" rules={[{ required: true, message: 'Please enter a description!' }]}>
                    <TextArea
                      maxLength={1000}
                      placeholder="Description of your subject"
                      style={{ minHeight: '240px', margin: '0.5rem 0' }}
                      size="large"
                    />
                  </Form.Item>
                </Form>
              </Modal>

              <Modal
                title="New Subject Description"
                visible={editSubDescVisible}
                onCancel={() => setEditSubDescVisible(false)}
                footer={[
                  <Button form="add-desc" key="submit" style={{ width: '100%' }} type="primary" htmlType="submit">
                    Add
                  </Button>,
                ]}
              >
                <Form onFinish={commitSubDescription} layout="vertical" name="add-desc" preserve={false}>
                  <Form.Item name="description" rules={[{ required: true, message: 'Please enter a description!' }]}>
                    <TextArea
                      maxLength={1000}
                      placeholder="Description of your subject"
                      style={{ minHeight: '240px', margin: '0.5rem 0' }}
                      size="large"
                    />
                  </Form.Item>
                </Form>
              </Modal>

              <Modal
                title="New Subject Price"
                visible={editSubPriceVisible}
                onCancel={() => setEditSubPriceVisible(false)}
                footer={[
                  <Button form="add-price" key="submit" style={{ width: '100%' }} type="primary" htmlType="submit">
                    Change
                  </Button>,
                ]}
              >
                <Form onFinish={commitSubPrice} layout="vertical" name="add-price" preserve={false}>
                  <Form.Item
                    name="price"
                    rules={[{ required: true, message: 'Please provide how much you wish your subject to cost.' }]}
                  >
                    <InputNumber
                      placeholder="Desired Cost of a Lesson"
                      size="large"
                      style={{ width: '100%' }}
                      min={1}
                    />
                  </Form.Item>
                </Form>
              </Modal>
            </>
          )}
        </Content>
      )}
      <RequestLessonModal
        onOk={() => {
          setRequestLessonVisible(false);
        }}
        onCancel={() => setRequestLessonVisible(false)}
        visible={requestLessonVisible}
        type={props.type}
        profile={profile}
      />
    </Typography>
  );
}
