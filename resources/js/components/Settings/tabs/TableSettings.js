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
        console.log(result);
        if (!result.destination) return;

        let visible = this.state.visible;
        let hidden = this.state.hidden;

        let from = result.source.droppableId == 'visibleColumns' ? visible.obj.value : hidden.obj.value;
        let to = result.destination.droppableId == 'visibleColumns' ? visible.obj.value : hidden.obj.value;

        let [reorderedItem] = from.splice(result.source.index, 1);
        to.splice(result.destination.index, 0, reorderedItem);

        this.setState({
            visible: visible,
            hidden: hidden
        });
    }

    save = () => {
        var url = 'api/settings/bulk?token=' + window.token;

        Axios.post(url, {
            data: [
                {
                    name: 'visible_columns',
                    value: this.state.visible.obj.value
                },
                {
                    name: 'hidden_columns',
                    value: this.state.hidden.obj.value
                }
            ],
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
