import React from 'react';
import ReactDOM from 'react-dom';
import 'antd/dist/antd.css';
import './index.css';
import './App.css';
import Layout from "antd/es/layout";
import CardList from "./CardList";

const { Header, Footer, Content } = Layout;

ReactDOM.render(
  <React.StrictMode>
      <Layout>
          <Header>Header</Header>
          <Content style={{ padding: '25px', maxWidth: '1200px', margin: 'auto'}}>
              <div className="site-layout-content">
                  <CardList>
                  </CardList>
              </div>
          </Content>
          <Footer style={{ textAlign: 'center' }}>
              Alice foo↑ foo↑
          </Footer>
      </Layout>
  </React.StrictMode>,
  document.getElementById('root')
);
