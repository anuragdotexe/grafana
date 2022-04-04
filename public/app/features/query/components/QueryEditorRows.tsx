// Libraries
import React, { PureComponent } from 'react';

// Types
import {
  CoreApp,
  DataQuery,
  DataSourceInstanceSettings,
  EventBusExtended,
  HistoryItem,
  PanelData,
} from '@grafana/data';
import { QueryEditorRow } from './QueryEditorRow';
import { DragStart, Droppable, DropResult } from 'react-beautiful-dnd';
import { getDataSourceSrv, reportInteraction } from '@grafana/runtime';
import { ScratchpadContext } from 'app/features/scratchpad';

interface Props {
  // The query configuration
  queries: DataQuery[];
  dsSettings: DataSourceInstanceSettings;

  // Query editing
  onQueriesChange: (queries: DataQuery[]) => void;
  onAddQuery: (query: DataQuery) => void;
  onRunQueries: () => void;
  id: string;

  // Query Response Data
  data: PanelData;

  // Misc
  app?: CoreApp;
  history?: Array<HistoryItem<DataQuery>>;
  eventBus?: EventBusExtended;
}

export class QueryEditorRows extends PureComponent<Props> {
  static contextType = ScratchpadContext;
  declare context: React.ContextType<typeof ScratchpadContext>;

  componentDidMount() {
    this.context?.register(this.props.id, {
      onDragStart: this.onDragStart,
      onDragEnd: this.onDragEnd,
      getItem: (index) => this.props.queries[index],
      onAddItem: (index, query) => {
        const queries = [...this.props.queries];
        queries.splice(index, 0, query);
        this.props.onQueriesChange(queries);
      },
      onRemoveItem: (index) => {
        this.props.onQueriesChange([...this.props.queries.slice(0, index), ...this.props.queries.slice(index + 1)]);
      },
    });
  }

  onRemoveQuery = (query: DataQuery) => {
    this.props.onQueriesChange(this.props.queries.filter((item) => item !== query));
  };

  onChangeQuery(query: DataQuery, index: number) {
    const { queries, onQueriesChange } = this.props;

    // update query in array
    onQueriesChange(
      queries.map((item, itemIndex) => {
        if (itemIndex === index) {
          return query;
        }
        return item;
      })
    );
  }

  onDataSourceChange(dataSource: DataSourceInstanceSettings, index: number) {
    const { queries, onQueriesChange } = this.props;

    onQueriesChange(
      queries.map((item, itemIndex) => {
        if (itemIndex !== index) {
          return item;
        }

        if (item.datasource) {
          const previous = getDataSourceSrv().getInstanceSettings(item.datasource);

          if (previous?.type === dataSource.type) {
            return {
              ...item,
              datasource: { uid: dataSource.uid },
            };
          }
        }

        return {
          refId: item.refId,
          hide: item.hide,
          datasource: { uid: dataSource.uid },
        };
      })
    );
  }

  onDragStart = (result: DragStart) => {
    const { queries, dsSettings } = this.props;

    reportInteraction('query_row_reorder_started', {
      startIndex: result.source.index,
      numberOfQueries: queries.length,
      datasourceType: dsSettings.type,
    });
  };

  onDragEnd = (result: DropResult) => {
    const { queries, onQueriesChange, dsSettings } = this.props;

    if (!result || !result.destination) {
      return;
    }

    const startIndex = result.source.index;
    const endIndex = result.destination.index;
    if (startIndex === endIndex) {
      reportInteraction('query_row_reorder_canceled', {
        startIndex,
        endIndex,
        numberOfQueries: queries.length,
        datasourceType: dsSettings.type,
      });
      return;
    }

    // const update = Array.from(queries);
    // const [removed] = update.splice(startIndex, 1);
    // update.splice(endIndex, 0, removed);
    // onQueriesChange(update);

    reportInteraction('query_row_reorder_ended', {
      startIndex,
      endIndex,
      numberOfQueries: queries.length,
      datasourceType: dsSettings.type,
    });
  };

  render() {
    const { dsSettings, data, queries, app, history, eventBus, id } = this.props;

    return (
      <Droppable droppableId={id} direction="vertical">
        {(provided) => {
          return (
            <div ref={provided.innerRef} {...provided.droppableProps}>
              {queries.map((query, index) => {
                const dataSourceSettings = getDataSourceSettings(query, dsSettings);
                const onChangeDataSourceSettings = dsSettings.meta.mixed
                  ? (settings: DataSourceInstanceSettings) => this.onDataSourceChange(settings, index)
                  : undefined;

                return (
                  <QueryEditorRow
                    id={`${id}-${query.refId}`}
                    index={index}
                    key={query.refId}
                    data={data}
                    query={query}
                    dataSource={dataSourceSettings}
                    onChangeDataSource={onChangeDataSourceSettings}
                    onChange={(query) => this.onChangeQuery(query, index)}
                    onRemoveQuery={this.onRemoveQuery}
                    onAddQuery={this.props.onAddQuery}
                    onRunQuery={this.props.onRunQueries}
                    queries={queries}
                    app={app}
                    history={history}
                    eventBus={eventBus}
                  />
                );
              })}
              {provided.placeholder}
            </div>
          );
        }}
      </Droppable>
    );
  }
}

const getDataSourceSettings = (
  query: DataQuery,
  groupSettings: DataSourceInstanceSettings
): DataSourceInstanceSettings => {
  if (!query.datasource) {
    return groupSettings;
  }
  const querySettings = getDataSourceSrv().getInstanceSettings(query.datasource);
  return querySettings || groupSettings;
};
