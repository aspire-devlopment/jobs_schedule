

Create Database [TaskScheduler]

USE [TaskScheduler]
GO

CREATE TABLE [dbo].[Jobs](
	[Id] [uniqueidentifier] NOT NULL,
	[Type] [nvarchar](100) NULL,
	[Payload] [nvarchar](max) NULL,
	[Priority] [nvarchar](10) NULL,
	[MaxRetries] [int] NULL,
	[RetryCount] [int] NULL,
	[Status] [nvarchar](50) NULL,
	[CreatedAt] [datetime] NULL,
	[UpdatedAt] [datetime] NULL,
PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]
GO

ALTER TABLE [dbo].[Jobs] ADD  DEFAULT ((0)) FOR [RetryCount]
GO

ALTER TABLE [dbo].[Jobs] ADD  DEFAULT (getdate()) FOR [CreatedAt]
GO


CREATE TABLE [dbo].[FailedJobs](
	[Id] [uniqueidentifier] NOT NULL,
	[OriginalJobId] [nvarchar](100) NULL,
	[Type] [nvarchar](100) NULL,
	[Payload] [nvarchar](max) NULL,
	[Priority] [nvarchar](10) NULL,
	[Reason] [nvarchar](max) NULL,
	[FailedAt] [datetime] NULL,
PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]
GO

ALTER TABLE [dbo].[FailedJobs] ADD  DEFAULT (getdate()) FOR [FailedAt]
GO
