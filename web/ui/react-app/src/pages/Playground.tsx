import React, { FC, useEffect, useState } from 'react';
import { RouteComponentProps } from '@reach/router';
import {
  Alert,
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
import { useDebouncedCallback } from 'use-debounce';
import ReactMarkdown from 'react-markdown';
import { Light as SyntaxHighlighter } from 'react-syntax-highlighter';
import markdown from 'react-syntax-highlighter/dist/esm/languages/hljs/markdown';
import { atomOneLight } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import classnames from 'classnames';
import './Playground.css';
import demoAlert from './PlaygroundDemoAlert.json';
import { useFetch } from '../utils/useFetch';

SyntaxHighlighter.registerLanguage('markdown', markdown);

type TemplateConfig = {
  name: string;
  text: string;
};

type TemplatesConfig = {
  templates: [TemplateConfig];
};

type RenderApiData = {
  markdown: string;
};

interface RenderApiResponse {
  status: string;
  error?: string;
  errorType?: string;
  data?: RenderApiData;
}

const demoAlertJson = JSON.stringify(demoAlert, null, 2);

const initialInputs = {
  text: `{{ template "ding.link.content" . }}`,
  demoAlertJSON: demoAlertJson,
};

const Playground: FC<RouteComponentProps> = () => {
  const [leftActiveTab, setLeftActiveTab] = useState('1');
  const [rightActiveTab, setRightActiveTab] = useState('1');
  const [inputs, setInputs] = useState(initialInputs);

  const startRenderMarkdown = async () => {
    try {
      const res = await fetch('/api/v1/status/templates/render', {
        method: 'POST',
        body: JSON.stringify(inputs),
      });

      const json = (await res.json()) as RenderApiResponse;
      if (res.ok) {
        setRenderError(false);

        if (json.data) {
          setMarkdown(json.data.markdown);
        }
      } else {
        setRenderError(true);
        console.info(`Error rendering template: ${json.error}`);
      }
    } catch (e) {
      setRenderError(true);
      console.info(`Unhandled error rendering template: ${e.toString()}`);
    }
  };
  const [delayedRender] = useDebouncedCallback(() => {
    startRenderMarkdown();
  }, 250);

  useEffect(() => {
    delayedRender();
  }, [delayedRender, inputs]);

  const [markdown, setMarkdown] = useState('');
  const [renderError, setRenderError] = useState(false);

  const { response: templateConfigResp } = useFetch<TemplatesConfig>('/api/v1/status/templates');
  let templates: TemplateConfig[] = [];
  if (templateConfigResp && templateConfigResp.data && templateConfigResp.data.templates) {
    templates = templateConfigResp.data.templates;
  }
  const [currentTemplate, setCurrentTemplate] = useState({ text: initialInputs.text });
  const handleTemplateChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    try {
      const idx = parseInt(event.target.value, 10);
      const tpl = templates[idx];
      setCurrentTemplate({ text: tpl.text });
    } catch (e) {
      setCurrentTemplate({ text: initialInputs.text });
    }
  };
  const loadTemplate = () => {
    const newState = {
      ...inputs,
      ...{
        text: currentTemplate.text,
      },
    };
    setInputs(newState);
  };

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
        <Col>
          <Alert color="danger" hidden={!renderError}>
            Unable to render template, check the console for details.
          </Alert>
        </Col>
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
                <FormGroup>
                  <Label className="form-label">Markdown Text:</Label>
                  <Input
                    type="textarea"
                    className="text-monospace"
                    style={{ height: '500px' }}
                    value={inputs.text}
                    onChange={evt => setInputs({ ...inputs, ...{ text: evt.target.value } })}
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
                    onChange={evt => setInputs({ ...inputs, ...{ demoAlertJSON: evt.target.value } })}
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
                        <ReactMarkdown
                          source={markdown}
                          className={'markdown-content'}
                          parserOptions={{
                            gfm: false,
                            commonmark: true,
                          }}
                        />
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
