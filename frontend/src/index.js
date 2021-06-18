import React from 'react';
import ReactDOM from 'react-dom';
import 'antd/dist/antd.css';
import './index.css';
import './App.css';
import {Layout, Menu} from "antd";
import CardList from "./CardList";
import Profile from "./Profile.js";
import Rank from "./Rank.js";
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
} from "react-router-dom";

const {Header, Footer, Content} = Layout;

ReactDOM.render(
    <React.StrictMode>
        <Router>
            <Layout style={{background: 'white'}}>
                <Header>
                    <Menu theme="dark" mode="horizontal" defaultSelectedKeys="1">
                        <Menu.Item key="1">
                            <Link to="/">直播</Link>
                        </Menu.Item>
                        <Menu.Item key="2">
                            <Link to="/rank">周榜</Link>
                        </Menu.Item>
                    </Menu>
                </Header>
                <Content style={{padding: '25px', maxWidth: '1200px', margin: 'auto', width: '100%'}}>
                    <div className="site-layout-content">
                        <Switch>
                            <Route path="/profile/:uid">
                                <Profile/>
                            </Route>
                            <Route path="/rank">
                                <Rank/>
                            </Route>
                            <Route path="/">
                                <CardList/>
                            </Route>
                        </Switch>
                    </div>
                </Content>
                <Footer style={{textAlign: 'center'}}>
                    Alice foo↑ foo↑
                </Footer>
            </Layout>
        </Router>
    </React.StrictMode>,
    document.getElementById('root')
);
