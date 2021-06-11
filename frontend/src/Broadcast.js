import React from 'react';
import {Card} from 'antd';
import {UserOutlined} from '@ant-design/icons';
import 'antd/dist/antd.css';
import './App.css';

const {Meta} = Card;

class Broadcast extends React.Component {
    render() {
        const item = this.props.item
        const modalClick = this.props.modalClick
        return (
            <Card
                cover={
                    <img
                        referrerPolicy="no-referrer"
                        src={item.usercover}
                        alt="cover"
                        onMouseOver={e => (e.currentTarget.src = item.keyframe)}
                        onMouseOut={e => (e.currentTarget.src = item.usercover)}
                        style={{width: 240, height: 125}}
                    />
                }
                hoverable={true}
                bordered={false}
                style={{width: 240, height: 220}}
                key={item.uid}
                onClick={(e) => modalClick(e, item)}
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
