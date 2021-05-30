import React from 'react';
import './App.css';
import 'antd/dist/antd.css';
import './index.css';
import 'antd/dist/antd.css';
import './index.css';
import {List, Card} from 'antd';
import {UserOutlined} from '@ant-design/icons';

const {Meta} = Card;

class CardList extends React.Component {
    state = {
        loading: true,
        data: [],
    };

    componentDidMount() {
        clearInterval(this.timer)
        this.getData(res => {
            this.setState({
                data: res,
                loading: false,
            })
        })
        this.timer = setInterval(() => {
            this.getData(res => {
                this.setState({
                    data: res,
                    loading: false,
                })
            })
        }, 30000)
    };

    componentWillUnmount() {
        this.timer && clearInterval(this.timer)
    }

    getData = callback => {
        fetch('https://livevup.com/api/online')
            .then(res => res.json())
            .then(res => callback(res.list))
    }

    render() {
        return (
            <List
                grid={{
                    gutter: 16,
                    xs: 1,
                    sm: 1,
                    md: 2,
                    lg: 4,
                    xl: 4,
                    xxl: 4,
                }}
                dataSource={this.state.data}
                loading={this.state.loading}
                renderItem={item => (
                    <List.Item>
                        <Card
                            cover={
                                <img
                                    src={item.usercover}
                                    alt="cover"
                                    referrerPolicy="no-referrer"
                                    onMouseOver={e => (e.currentTarget.src = item.keyframe)}
                                    onMouseOut={e => (e.currentTarget.src = item.usercover)}
                                />
                            }
                            hoverable={true}
                            style={{width: 240}}
                        >
                            <Meta
                                title={<a
                                    href={"https://live.bilibili.com/" + item.roomid}
                                    target="_blank"
                                    rel="noreferrer"
                                >{item.title}</a>}
                                description={
                                    <div>
                                        <a
                                            href={"https://space.bilibili.com/" + item.uid}
                                            target="_blank"
                                            rel="noreferrer"
                                            style={{display: 'inline-block', maxWidth: '75%'}}
                                        >
                                            {item.uname}
                                        </a>
                                        <span style={{float: 'right'}}>
                                      <UserOutlined/>
                                            {item.participantDuring10Min}
                                  </span>
                                    </div>
                                }
                            />
                        </Card>
                    </List.Item>
                )}
            />
        );
    }
}

export default CardList;
