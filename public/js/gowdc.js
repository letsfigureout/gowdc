(function() {
    // Create the connector object
    var myConnector = tableau.makeConnector();

    // Define the schema
    myConnector.getSchema = function(schemaCallback) {
        var cols = [{
            id: "id",
            dataType: tableau.dataTypeEnum.string
        }, {
            id: "Date",
            alias: "Date",
            dataType: tableau.dataTypeEnum.date
        }, {
            id: "Open",
            alias: "Open",
            dataType: tableau.dataTypeEnum.float
        }, {
            id: "High",
            alias: "High",
            dataType: tableau.dataTypeEnum.float
        }, {
            id: "Low",
            alias: "Low",
            dataType: tableau.dataTypeEnum.float
        }, {
            id: "Close",
            alias: "Cloe",
            dataType: tableau.dataTypeEnum.float
        }, {
            id: "AdjClose",
            alias: "AdjClose",
            dataType: tableau.dataTypeEnum.float
        }, {
            id: "Volume",
            alias: "Volume",
            dataType: tableau.dataTypeEnum.int
        }];

        var tableSchema = {
            id: "stockdata",
            alias: "Historical stock data",
            columns: cols
        };

        schemaCallback([tableSchema]);
    };

    // Download the data
    myConnector.getData = function(table, doneCallback) {
        var tickerList = JSON.parse(tableau.connectionData)
        last = tickerList.length
        var param = ""

        tickerList.forEach(function (tkr, index) {
            if (index == 0) {
                param = param + "?ticker=" + tkr + "&";
            } else if (index == tickerList.length) {
                param = param + "ticker=" + tkr;
            } else {
                param = param + "ticker=" + tkr + "&";
            }

        })

        $.getJSON("http://" + window.location.host + "/api" + param, function(resp) {
            tableData = [];

            // Iterate over the JSON object
            for (var tickerName in resp.Stock) {
                data = resp.Stock[tickerName]
                for (var xi=0; xi < data.PriceHistory.length; xi++) {
                    var row = data.PriceHistory[xi]
                    tableData.push({
                        "id": tickerName,
                        "Date":     row.Date,
                        "Open":     row.Open,
                        "High":     row.High,
                        "Low":      row.Low,
                        "Close":    row.Close,
                        "AdjClose": row.AdjClose,
                        "Volume":   row.Volume,
                    });
                }
                table.appendRows(tableData);

            }
            doneCallback();
        });
    };

    tableau.registerConnector(myConnector);

    // Create event listeners for when the user submits the form
    $(document).ready(function() {
        $("#submitButton").click(function(e) {
            e.preventDefault()

            var tickerObj = $('#tickers').val().trim();

            if (tickerObj != undefined && tickerObj.length > 0) {
                tickerObj = tickerObj.split(',')
                tableau.connectionName = "Stock Data"; // This will be the data source name in Tableau
                tableau.connectionData = JSON.stringify(tickerObj);
                tableau.submit(); // This sends the connector object to Tableau
            } else {
                $('#errorMsg').html('ticker name cannot be empty');
            }
        });
    });
})();

