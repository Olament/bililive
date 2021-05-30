import React from 'react';
import {Card} from 'antd';
import {UserOutlined} from '@ant-design/icons';
import 'antd/dist/antd.css';
import './App.css';

const {Meta} = Card;

class Broadcast extends React.Component {
    render() {
        const item = this.props.item
        return (
            <Card
                cover={
                    <img
                        src={item.usercover}
                        alt="cover"
                        referrerPolicy="no-referrer"
                        onMouseOver={e => (e.currentTarget.src = item.keyframe)}
                        onMouseOut={e => (e.currentTarget.src = item.usercover)}
                        style={{width: 240, height: 125}}
                    />
                }
                hoverable={true}
                style={{width: 240, height: 220}}
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
        )
    }
}

export default Broadcast;
