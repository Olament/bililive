import React from 'react';
import {Card, Tooltip} from 'antd';
import {Link} from "react-router-dom";
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
                bordered={true}
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
                            <Link
                                to={'/profile/'+item.uid}
                                rel="noreferrer"
                                style={{display: 'inline-block', maxWidth: '75%'}}>
                                {item.uname}
                            </Link>
                            <span style={{float: 'right'}}>
                                <Tooltip title="十分钟互动人数">
                                    <UserOutlined/>
                                    {item.participantDuring10Min}
                                </Tooltip>
                            </span>
                        </div>
                    }
                />
            </Card>
        )
    }
}

export default Broadcast;
