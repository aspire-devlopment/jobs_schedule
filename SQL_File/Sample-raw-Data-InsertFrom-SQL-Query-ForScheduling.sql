INSERT INTO [dbo].[Jobs] (
    [Id],
    [Type],
    [Payload],
    [Priority],
    [MaxRetries],
    [RetryCount],
    [Status],
    [CreatedAt],
    [UpdatedAt]
)
VALUES 
(NEWID(), 'send_email', '{"to":"user1@example.com","subject":"Welcome"}', 'HIGH', 3, 0, 'queued', GETDATE(), NULL),

(NEWID(), 'export_user_data', '{"userId":42,"format":"PDF"}', 'MEDIUM', 2, 1, 'queued', GETDATE(), NULL),

(NEWID(), 'process_payment', '{"orderId":12345,"amount":99.99}', 'HIGH', 3, 2, 'queued', GETDATE(), NULL),

(NEWID(), 'Notification', '{"to":"user5","message":"Daily digest"}', 'LOW', 1, 0, 'queued', GETDATE(), NULL),

(NEWID(), 'ReportGeneration', '{"reportType":"monthly"}', 'MEDIUM', 2, 2, 'queued', GETDATE(), NULL);
