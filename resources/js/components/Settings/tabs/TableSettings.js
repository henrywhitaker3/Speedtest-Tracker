import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Modal, Button, Tab } from 'react-bootstrap';
import Axios from 'axios';
import { toast } from 'react-toastify';
import { DragDropContext, Draggable, Droppable } from 'react-beautiful-dnd';

export default class TableSettings extends Component {
    constructor(props) {
        super(props)

        this.state = {
            visible: this.props.data[0],
            hidden: this.props.data[1],
        }
    }

    handleOnDragEnd = (result) => {
        if (!result.destination) return;

        let data = this.state.visible;
        let array = data.obj.value;

        let reorderedItem = array.splice(result.source.index, 1);
        array.splice(result.destination.index, 0, reorderedItem[0]);
        data.obj.value = array;

        console.log(array);

        this.setState({
            visible: data
        });
    }

    save = () => {
        var url = 'api/settings?token=' + window.token;

        Axios.post(url, {
            name: 'visible_columns',
            value: this.state.data.obj.value
        })
        .then((resp) => {
            console.log(resp);
        })
        .catch((err) => {
            console.log(err);
        })
    }

    render() {
        let visible = this.state.visible;
        let hidden = this.state.hidden;

        return (
            <Tab.Content>
                <div>
                    <p>{visible.obj.description}</p>

                    <DragDropContext onDragEnd={this.handleOnDragEnd}>
                        <Droppable droppableId="visibleColumns">
                            {(provided) => (
                                <ul className="visibleColumns" {...provided.droppableProps} ref={provided.innerRef}>
                                    {visible.obj.value.map((e, i) => {
                                        return (
                                            <Draggable draggableId={e} index={i} key={e}>
                                                {(provided) => (
                                                    <li key={e} ref={provided.innerRef} {...provided.draggableProps} {...provided.dragHandleProps}>{e}</li>
                                                )}
                                            </Draggable>
                                        );
                                    })}
                                    {provided.placeholder}
                                </ul>
                            )}
                        </Droppable>
                        <Droppable droppableId="hiddenColumns">
                            {(provided) => (
                                <ul className="hiddenColumns" {...provided.droppableProps} ref={provided.innerRef}>
                                    {hidden.obj.value.map((e, i) => {
                                        return (
                                            <Draggable draggableId={e} index={i} key={e}>
                                                {(provided) => (
                                                    <li key={e} ref={provided.innerRef} {...provided.draggableProps} {...provided.dragHandleProps}>{e}</li>
                                                )}
                                            </Draggable>
                                        );
                                    })}
                                    {provided.placeholder}
                                </ul>
                            )}
                        </Droppable>
                    </DragDropContext>

                    <div className="mt-3">
                        <button className="btn btn-primary" onClick={() => { this.save() }}>Save</button>
                    </div>
                </div>
            </Tab.Content>
        );
    }
}

if (document.getElementById('TableSettings')) {
    ReactDOM.render(<TableSettings />, document.getElementById('TableSettings'));
}
