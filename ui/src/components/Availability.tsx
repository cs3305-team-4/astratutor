import { Table } from "antd"


export interface AvailabilityProps {
    hours: boolean[]
    editable: boolean
    onUpdate(hours: boolean[]): void
}

interface availabilityRow {
  mon: boolean
  tue: boolean
  wed: boolean
  thu: boolean
  fri: boolean
  sat: boolean
  sun: boolean
}

function flatHoursToTable(hoursIn: boolean[]) {
  let hours = [...hoursIn]

  for(let i = 0; i < (24*7); i++) {
    if (i > (hours.length-1)) {
      hours[i] = false
    }
  }
  
  let ret = []
  let k = 0
  for(let i = 0; i < 7; i++) {
    ret.push({
      key: k.toString(),
      mon: hours[i],
      tue: hours[24*1+i],
      wed: hours[24*2+i],
      thu: hours[24*3+i],
      fri: hours[24*4+i],
      sat: hours[24*5+i],
      sun: hours[24*6+i]
    })

    k++
  }

  return ret
}

export default function Availability() {
  const columns = [
    {
      title: 'Monday',
      dataIndex: '0',
      key: '0',
    },
    {
      title: 'Tuesday',
      dataIndex: '1',
      key: '1',
    },
    {
      title: 'Wednesday',
      dataIndex: '2',
      key: '2',
    },
    {
      title: 'Thursday',
      dataIndex: '3',
      key: '3',
    },
    {
      title: 'Friday',
      dataIndex: '4',
      key: '4',
    },
    {
      title: 'Saturday',
      dataIndex: '5',
      key: '5',
    },
    {
      title: 'Sunday',
      dataIndex: '6',
      key: '6',
    },
  ]



  return (
    <Table>

    </Table>
    <Layout>
      <Typography>
        <Title>
          Profile
        </Title>
      </Typography>
    </Layout>
  )
}