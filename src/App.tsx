import React, { useState } from 'react';
import { Button, Layout, Row, Col, Typography, Space, Input } from 'antd';
import './App.css';

const { Header, Footer, Sider, Content } = Layout;
const { Title, Text } = Typography;

function App() {

  const [query, setQuery] = useState('q=avenger&cat=207');

  return (
    <Layout
      style={{
        height: '100vh',
        textAlign: 'center',
      }}
    >
      <Content>

        <Space direction='vertical' size={100}>
          <Title level={4}>1. Search anything you like in <a href='https://thepiratebay.org'>https://thepiratebay.org</a></Title>

          <div>
            <Title level={4}>2. For example, search with 'avenger' in 'HD - Movies' category</Title>
            <Text><a href='https://thepiratebay.org/search.php?q=avenger&cat=207'>https://thepiratebay.org/search.php?q=avenger&cat=207</a></Text>
          </div>

          <div>
            <Title level={4}>3. Copy and past parameters in URL</Title>
            <Row justify='center'>
              <Col
                style={{
                  margin: 'auto 0%'
                }}
              >
                <Text>https://thepiratebay.org/search.php?</Text>
              </Col>
              <Col>
                <Input placeholder={query} onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                  setQuery(e.target.value)
                }} />
              </Col>
            </Row>
          </div>

          <div>
            <Title level={4}>4. Get RSS link</Title>
            <Row justify='center'>
              <Col
                style={{
                  margin: 'auto 0%'
                }}
              >
                <Button type='link'>https://todo.com?{query}</Button>
              </Col>
            </Row>
          </div>

        </Space>

      </Content>
    </Layout>
  );
}

export default App;
