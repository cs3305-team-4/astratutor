import React, { useContext } from 'react';
import { Button, Switch, Table, Modal } from 'antd';

import { CheckOutlined, StopOutlined } from '@ant-design/icons';
import { AccountType, ProfileResponseDTO } from '../api/definitions';
import { useAsync } from 'react-async-hook';
import { APIContext } from '../api/api';

export interface AvailabilityProps {
  availability: boolean[];
  hideUnavailable: boolean;
  editable?: boolean;
  onChange?(hours: boolean[]): void;
}

interface AvailabilitySlotProps {
  available: boolean;
  editable: boolean;
  onChange(available: boolean): void;
}

const AvailabilitySlot: React.FC<AvailabilitySlotProps> = (props: AvailabilitySlotProps) => {
  if (props.editable) {
    return (
      <Switch
        style={{ margin: '0 auto' }}
        size="small"
        checked={props.available}
        onChange={(checked: boolean) => {
          console.log(checked);
          props.onChange(checked);
        }}
      />
    );
  } else {
    if (props.available) {
      return <CheckOutlined style={{ color: 'green' }} />;
    } else {
      return <StopOutlined style={{ color: 'red' }} />;
    }
    return <h1>{props.available === true ? 'available' : 'not available'}</h1>;
  }
};

// Takes in an array of hours and pads it to meet the length of 168 (24*7)
function padHours(hoursOld: boolean[]) {
  const hours = [...hoursOld];

  // fill in any hours that were missing
  for (let i = 0; i < 24 * 7; i++) {
    if (i > hours.length - 1) {
      hours[i] = false;
    }
  }

  return hours;
}

export const Availability: React.FC<AvailabilityProps> = (props: AvailabilityProps) => {
  const availability = padHours(props.availability);

  const columns = [
    {
      title: 'Time',
      dataIndex: 'time',
      key: 'time',
    },
    {
      title: 'Mon',
      dataIndex: 'mon',
      key: 'mon',
    },
    {
      title: 'Tue',
      dataIndex: 'tue',
      key: '1',
    },
    {
      title: 'Wed',
      dataIndex: 'wed',
      key: '2',
    },
    {
      title: 'Thu',
      dataIndex: 'thu',
      key: '3',
    },
    {
      title: 'Fri',
      dataIndex: 'fri',
      key: '4',
    },
    {
      title: 'Sat',
      dataIndex: 'sat',
      key: '5',
    },
    {
      title: 'Sun',
      dataIndex: 'sun',
      key: '6',
    },
  ];

  // set an hour to be available, merge with existing availability and calls onChange
  const onChangeAvailSingleHour = async (hour: number, available: boolean) => {
    const newHours = [...availability];
    newHours[hour] = available;
    props.onChange(newHours);
  };

  // convert flat hours table into table rows
  const rows = [];
  let k = 0;
  for (let i = 8; i < 22; i++) {
    const row = {
      key: k.toString(),
      time: `${i.toString().padStart(2, '0')}:00`,
    };

    let rowHasAvailable = false;
    const days = ['mon', 'tue', 'wed', 'thu', 'fri', 'sat', 'sun'];
    days.forEach((day: string, day_i: number) => {
      const hour_i = 24 * day_i + i;

      if (availability[hour_i] === true) {
        rowHasAvailable = true;
      }

      row[day] = (
        <AvailabilitySlot
          editable={props.editable}
          onChange={(available: boolean) => {
            onChangeAvailSingleHour(hour_i, available);
          }}
          available={availability[hour_i]}
        />
      );
    });

    if (props.hideUnavailable && !props.editable) {
      if (rowHasAvailable) {
        rows.push(row);
      }
    } else {
      rows.push(row);
    }

    k++;
  }

  return (
    <Table
      style={{ marginBottom: '0.5rem' }}
      locale={{
        emptyText: 'No times available',
      }}
      size={'small'}
      pagination={false}
      columns={columns}
      dataSource={rows}
    ></Table>
  );
};
