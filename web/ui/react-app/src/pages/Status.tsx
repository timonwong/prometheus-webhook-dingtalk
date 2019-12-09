import React, { FC, Fragment } from 'react';
import { RouteComponentProps } from '@reach/router';
import { Table } from 'reactstrap';
import { withStatusIndicator } from '../withStatusIndicator';
import { useFetch } from '../utils/useFetch';

const sectionTitles = ['Runtime Information', 'Build Information'];

interface StatusConfig {
  [k: string]: { title?: string; customizeValue?: (v: any, key: string) => any; customRow?: boolean; skip?: boolean };
}

type StatusPageState = { [k: string]: string };

interface StatusPageProps {
  data?: StatusPageState[];
}

export const statusConfig: StatusConfig = {
  startTime: { title: 'Start time', customizeValue: (v: string) => new Date(v).toUTCString() },
  CWD: { title: 'Working directory' },
  reloadConfigSuccess: {
    title: 'Configuration reload',
    customizeValue: (v: boolean) => (v ? 'Successful' : 'Unsuccessful'),
  },
  lastConfigTime: { title: 'Last successful configuration reload' },
  goroutineCount: { title: 'Goroutines' },
};

export const StatusContent: FC<StatusPageProps> = ({ data = [] }) => {
  return (
    <>
      {data.map((statuses, i) => {
        return (
          <Fragment key={i}>
            <h2>{sectionTitles[i]}</h2>
            <Table className="h-auto" size="sm" bordered striped>
              <tbody>
                {Object.entries(statuses).map(([k, v], i) => {
                  const { title = k, customizeValue = (val: any) => val, customRow, skip } = statusConfig[k] || {};
                  if (skip) {
                    return null;
                  }
                  if (customRow) {
                    return customizeValue(v, k);
                  }
                  return (
                    <tr key={k}>
                      <th className="capitalize-title" style={{ width: '35%' }}>
                        {title}
                      </th>
                      <td className="text-break">{customizeValue(v, title)}</td>
                    </tr>
                  );
                })}
              </tbody>
            </Table>
          </Fragment>
        );
      })}
    </>
  );
};
const StatusWithStatusIndicator = withStatusIndicator(StatusContent);

StatusContent.displayName = 'Status';

const Status: FC<RouteComponentProps> = () => {
  const path = '/api/v1';
  const status = useFetch<StatusPageState>(`${path}/status/runtimeinfo`);
  const runtime = useFetch<StatusPageState>(`${path}/status/buildinfo`);

  let data;
  if (status.response.data && runtime.response.data) {
    data = [status.response.data, runtime.response.data];
  }

  return (
    <StatusWithStatusIndicator
      data={data}
      isLoading={status.isLoading || runtime.isLoading}
      error={status.error || runtime.error}
    />
  );
};

export default Status;
