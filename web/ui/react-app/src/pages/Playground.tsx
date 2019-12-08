import React, { FC, useRef, useState } from 'react';
import { RouteComponentProps } from '@reach/router';
import {
  Button,
  Col,
  Form,
  FormGroup,
  Input,
  InputGroup,
  InputGroupAddon,
  Label,
  Nav,
  NavItem,
  NavLink,
  Row,
  TabContent,
  TabPane,
} from 'reactstrap';
import _ from 'lodash';
import ReactMarkdown from 'react-markdown';
import { Light as SyntaxHighlighter } from 'react-syntax-highlighter';
import markdown from 'react-syntax-highlighter/dist/esm/languages/hljs/markdown';
import { atomOneLight } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import classnames from 'classnames';
import './Playground.css';
import { useFetch } from '../utils/useFetch';

SyntaxHighlighter.registerLanguage('markdown', markdown);

type TemplateConfig = {
  name: string;
  title: string;
  text: string;
};

type TemplatesConfig = {
  templates: [TemplateConfig];
};

type APIRenderData = {
  markdown: string;
};

interface APIRenderResponse {
  status: string;
  data?: APIRenderData;
}

const demoAlertJSON = `{
    "receiver": "admins",
    "status": "firing",
    "alerts": [
        {
            "status": "firing",
            "labels": {
                "alertname": "something_happened",
                "env": "prod",
                "instance": "server01.int:9100",
                "job": "node",
                "service": "prometheus_bot",
                "severity": "warning",
                "supervisor": "runit"
            },
            "annotations": {
                "summary": "Oops, something happened!"
            },
            "startsAt": "2016-04-27T20:46:37.903Z",
            "endsAt": "0001-01-01T00:00:00Z",
            "generatorURL": "https://example.com/graph#..."
        },
        {
            "status": "firing",
            "labels": {
                "alertname": "something_happened",
                "env": "staging",
                "instance": "server02.int:9100",
                "job": "node",
                "service": "prometheus_bot",
                "severity": "warning",
                "supervisor": "runit"
            },
            "annotations": {
                "summary": "Oops, something happend!"
            },
            "startsAt": "2016-04-27T20:49:37.903Z",
            "endsAt": "0001-01-01T00:00:00Z",
            "generatorURL": "https://example.com/graph#..."
        }
    ],
    "groupLabels": {
        "alertname": "something_happened",
        "instance": "server01.int:9100"
    },
    "commonLabels": {
        "alertname": "something_happened",
        "job": "node",
        "service": "prometheus_bot",
        "severity": "warning",
        "supervisor": "runit"
    },
    "commonAnnotations": {
        "summary": "runit service prometheus_bot restarted, server01.int:9100"
    },
    "externalURL": "https://alert-manager.example.com",
    "version": "3"
}`;

const Playground: FC<RouteComponentProps> = () => {
  const [leftActiveTab, setLeftActiveTab] = useState('1');
  const [rightActiveTab, setRightActiveTab] = useState('1');
  const [inputs, setInputs] = useState({
    title: `{{ template "ding.link.title" . }}`,
    text: `{{ template "ding.link.content" . }}`,
    demoAlertJSON: demoAlertJSON,
  });

  const delayedRender = useRef(_.debounce(() => sendDelayedRender(), 1000)).current;
  const [markdown, setMarkdown] = useState('');
  const sendDelayedRender = async () => {
    const res = await fetch('/api/v1/status/templates/render', {
      method: 'POST',
      body: JSON.stringify(inputs),
    });
    if (res.ok) {
      const json = (await res.json()) as APIRenderResponse;
      if (json.data) {
        setMarkdown(json.data.markdown);
      }
    }
  };

  const { templateConfigResp } = useFetch<TemplatesConfig>('/api/v1/status/templates');

  let templateValue: string;
  let templates: TemplateConfig[] = [];
  if (templateConfigResp && templateConfigResp.data && templateConfigResp.data.templates) {
    templates = templateConfigResp.data.templates;
    templateValue = '0';
  }

  const handleTemplateChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    templateValue = event.target.value;
  };
  const loadTemplate = () => {
    if (templateValue) {
      const tpl = templates[parseInt(templateValue)];
      const newState = {
        ...inputs,
        ...{
          title: tpl.title,
          text: tpl.text,
        },
      };
      setInputs(newState);
    }
  };

  sendDelayedRender();

  const templateLoaderComponent = (
    <InputGroup>
      <InputGroupAddon addonType="prepend">Templates:</InputGroupAddon>
      <Input type="select" onChange={handleTemplateChange}>
        {templates.map((c, id) => (
          <option key={id} value={id}>
            {c.name}
          </option>
        ))}
      </Input>
      <InputGroupAddon addonType="append">
        <Button color="primary" onClick={loadTemplate}>
          Load
        </Button>
      </InputGroupAddon>
    </InputGroup>
  );

  return (
    <div className="panel">
      <Row>
        <Col>
          <Form>
            <FormGroup>{templateLoaderComponent}</FormGroup>
          </Form>
        </Col>
        <Col />
      </Row>
      <Row>
        <Col>
          <div className="preview-nav">
            <Nav pills>
              <NavItem>
                <NavLink
                  className={classnames({ active: leftActiveTab === '1' })}
                  onClick={() => setLeftActiveTab('1')}
                  href="#"
                >
                  Template
                </NavLink>
              </NavItem>
              <NavItem>
                <NavLink
                  className={classnames({ active: leftActiveTab === '2' })}
                  onClick={() => setLeftActiveTab('2')}
                  href="#"
                >
                  Alert JSON
                </NavLink>
              </NavItem>
            </Nav>
          </div>
          <TabContent activeTab={leftActiveTab}>
            <TabPane tabId="1">
              <Form>
                {/*<FormGroup>*/}
                {/*  <Label className="form-label">Markdown Title:</Label>*/}
                {/*  <Input*/}
                {/*    type="textarea"*/}
                {/*    className="text-monospace"*/}
                {/*    value={inputs.title}*/}
                {/*    onChange={evt => {*/}
                {/*      setInputs({ ...inputs, ...{ title: evt.target.value } });*/}
                {/*      delayedRender();*/}
                {/*    }}*/}
                {/*  />*/}
                {/*</FormGroup>*/}
                <FormGroup>
                  <Label className="form-label">Markdown Text:</Label>
                  <Input
                    type="textarea"
                    className="text-monospace"
                    style={{ height: '500px' }}
                    value={inputs.text}
                    onChange={evt => {
                      setInputs({ ...inputs, ...{ text: evt.target.value } });
                      delayedRender();
                    }}
                  />
                </FormGroup>
              </Form>
            </TabPane>
            <TabPane tabId="2">
              <Form>
                <FormGroup>
                  <Label className="form-label">Demo Prometheus Alert (JSON):</Label>
                  <Input
                    type="textarea"
                    className="text-monospace"
                    style={{ height: '500px' }}
                    value={inputs.demoAlertJSON}
                    onChange={evt => {
                      setInputs({ ...inputs, ...{ demoAlertJSON: evt.target.value } });
                      delayedRender();
                    }}
                  />
                </FormGroup>
              </Form>
            </TabPane>
          </TabContent>
        </Col>
        <Col md={6}>
          <div className="preview-nav">
            <Nav pills>
              <NavItem>
                <NavLink
                  className={classnames({ active: rightActiveTab === '1' })}
                  onClick={() => setRightActiveTab('1')}
                  href="#"
                >
                  Preview
                </NavLink>
              </NavItem>
              <NavItem>
                <NavLink
                  className={classnames({ active: rightActiveTab === '2' })}
                  onClick={() => setRightActiveTab('2')}
                  href="#"
                >
                  Markdown
                </NavLink>
              </NavItem>
            </Nav>
          </div>
          <TabContent activeTab={rightActiveTab}>
            <TabPane tabId="1">
              <Form>
                <FormGroup>
                  <Label className="form-label">Preview:</Label>
                  <div className="preview">
                    <div className="clearfix preview-content-area">
                      <div className="message-bubble">
                        <ReactMarkdown source={markdown} className={'markdown-content'} />
                      </div>
                    </div>
                  </div>
                </FormGroup>
              </Form>
            </TabPane>
            <TabPane tabId="2">
              <Form>
                <FormGroup>
                  <Label className="form-label">Markdown:</Label>
                  <div className="preview-markdown">
                    <SyntaxHighlighter language="markdown" style={atomOneLight} customStyle={{ height: '500px' }}>
                      {markdown}
                    </SyntaxHighlighter>
                  </div>
                </FormGroup>
              </Form>
            </TabPane>
          </TabContent>
        </Col>
      </Row>
    </div>
  );
};

export default Playground;
