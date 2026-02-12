export type TagOption = {
  value: string
  labelKey: string
  parent?: string // родительский тег для иерархии
  category?: string // основная категория (IT, GameDev, Startup и т.д.)
}

export type TagGroup = {
  key: string
  labelKey: string
  options: TagOption[]
  priority?: number // приоритет отображения (1 - высший)
  category?: 'primary' | 'secondary' // основные или дополнительные теги
}

// Иерархия тегов: category > subcategory > specialization > detail (4 уровня)
export type TagHierarchy = {
  value: string
  labelKey: string
  children?: TagHierarchy[]
  searchTerms?: string[] // дополнительные термины для поиска
}

// Иерархическая таксономия тегов
export const tagTaxonomy: TagHierarchy[] = [
  // ===== IT =====
  {
    value: 'IT',
    labelKey: 'tags.categories.it',
    searchTerms: ['технологии', 'разработка', 'tech'],
    children: [
      {
        value: 'Backend',
        labelKey: 'tags.specializations.backend',
        children: [
          { value: 'Backend+Node.js', labelKey: 'tags.details.backendNodejs' },
          { value: 'Backend+Python', labelKey: 'tags.details.backendPython' },
          { value: 'Backend+Go', labelKey: 'tags.details.backendGo' },
          { value: 'Backend+C#', labelKey: 'tags.details.backendCsharp' },
          { value: 'Backend+Java', labelKey: 'tags.details.backendJava' },
          { value: 'Backend+PHP', labelKey: 'tags.details.backendPhp' },
          { value: 'Backend+Ruby', labelKey: 'tags.details.backendRuby' }
        ]
      },
      {
        value: 'Frontend',
        labelKey: 'tags.specializations.frontend',
        children: [
          { value: 'Frontend+React', labelKey: 'tags.details.frontendReact' },
          { value: 'Frontend+Vue', labelKey: 'tags.details.frontendVue' },
          { value: 'Frontend+Angular', labelKey: 'tags.details.frontendAngular' },
          { value: 'Frontend+Next.js', labelKey: 'tags.details.frontendNextjs' },
          { value: 'Frontend+Nuxt', labelKey: 'tags.details.frontendNuxt' }
        ]
      },
      {
        value: 'Fullstack',
        labelKey: 'tags.specializations.fullstack',
        children: [
          { value: 'Fullstack+MERN', labelKey: 'tags.details.fullstackMern' },
          { value: 'Fullstack+MEAN', labelKey: 'tags.details.fullstackMean' },
          { value: 'Fullstack+Django', labelKey: 'tags.details.fullstackDjango' }
        ]
      },
      {
        value: 'Mobile',
        labelKey: 'tags.specializations.mobile',
        children: [
          { value: 'Mobile+iOS', labelKey: 'tags.details.mobileIos' },
          { value: 'Mobile+Android', labelKey: 'tags.details.mobileAndroid' },
          { value: 'Mobile+React Native', labelKey: 'tags.details.mobileReactNative' },
          { value: 'Mobile+Flutter', labelKey: 'tags.details.mobileFlutter' }
        ]
      },
      {
        value: 'DevOps',
        labelKey: 'tags.specializations.devops',
        children: [
          { value: 'DevOps+AWS', labelKey: 'tags.details.devopsAws' },
          { value: 'DevOps+Azure', labelKey: 'tags.details.devopsAzure' },
          { value: 'DevOps+GCP', labelKey: 'tags.details.devopsGcp' },
          { value: 'DevOps+Kubernetes', labelKey: 'tags.details.devopsK8s' },
          { value: 'DevOps+Docker', labelKey: 'tags.details.devopsDocker' }
        ]
      },
      {
        value: 'Data',
        labelKey: 'tags.specializations.data',
        children: [
          { value: 'Data+Analytics', labelKey: 'tags.details.dataAnalytics' },
          { value: 'Data+Science', labelKey: 'tags.details.dataScience' },
          { value: 'Data+Engineering', labelKey: 'tags.details.dataEngineering' },
          { value: 'Data+ML/AI', labelKey: 'tags.details.dataMlAi' }
        ]
      },
      {
        value: 'QA',
        labelKey: 'tags.specializations.qa',
        children: [
          { value: 'QA+Manual', labelKey: 'tags.details.qaManual' },
          { value: 'QA+Automation', labelKey: 'tags.details.qaAutomation' },
          { value: 'QA+Performance', labelKey: 'tags.details.qaPerformance' }
        ]
      },
      {
        value: 'Security',
        labelKey: 'tags.specializations.security',
        children: [
          { value: 'Security+InfoSec', labelKey: 'tags.details.securityInfosec' },
          { value: 'Security+AppSec', labelKey: 'tags.details.securityAppsec' },
          { value: 'Security+DevSecOps', labelKey: 'tags.details.securityDevsecops' }
        ]
      },
      {
        value: 'Design',
        labelKey: 'tags.specializations.design',
        children: [
          { value: 'Design+UX', labelKey: 'tags.details.designUx' },
          { value: 'Design+UI', labelKey: 'tags.details.designUi' },
          { value: 'Design+Product', labelKey: 'tags.details.designProduct' }
        ]
      }
    ]
  },

  // ===== GAMEDEV =====
  {
    value: 'GameDev',
    labelKey: 'tags.categories.gamedev',
    searchTerms: ['игры', 'games', 'геймдев'],
    children: [
      {
        value: 'Unity',
        labelKey: 'tags.specializations.unity',
        children: [
          { value: 'Unity+C#', labelKey: 'tags.details.unityCsharp' },
          { value: 'Unity+Mobile', labelKey: 'tags.details.unityMobile' },
          { value: 'Unity+PC', labelKey: 'tags.details.unityPc' }
        ]
      },
      {
        value: 'Unreal Engine',
        labelKey: 'tags.specializations.unreal',
        children: [
          { value: 'UE5+C++', labelKey: 'tags.details.ue5Cpp' },
          { value: 'UE5+C#', labelKey: 'tags.details.ue5Csharp' },
          { value: 'UE5+Blueprints', labelKey: 'tags.details.ue5Blueprints' }
        ]
      },
      {
        value: 'Game Design',
        labelKey: 'tags.specializations.gameDesign',
        children: [
          { value: 'Game Design+Level', labelKey: 'tags.details.levelDesign' },
          { value: 'Game Design+Narrative', labelKey: 'tags.details.narrativeDesign' },
          { value: 'Game Design+Systems', labelKey: 'tags.details.systemsDesign' }
        ]
      },
      {
        value: 'Game Art',
        labelKey: 'tags.specializations.gameArt',
        children: [
          {
            value: 'Game Art+3D',
            labelKey: 'tags.details.gameArt3d',
            children: [
              { value: 'Game Art+3D+Modeling+Character', labelKey: 'tags.details.gameArt3dCharacterArtist' },
              { value: 'Game Art+3D+Modeling+Environment', labelKey: 'tags.details.gameArt3dEnvironmentArtist' },
              { value: 'Game Art+3D+Modeling+Props', labelKey: 'tags.details.gameArt3dPropsArtist' }
            ]
          },
          {
            value: 'Game Art+2D',
            labelKey: 'tags.details.gameArt2d',
            children: [
              { value: 'Game Art+2D+Concept+Character', labelKey: 'tags.details.gameArt2dConceptCharacter' },
              { value: 'Game Art+2D+Concept+Environment', labelKey: 'tags.details.gameArt2dConceptEnvironment' },
              { value: 'Game Art+2D+Illustration', labelKey: 'tags.details.gameArt2dIllustration' }
            ]
          },
          {
            value: 'Game Art+Animation',
            labelKey: 'tags.details.gameArtAnimation',
            children: [
              { value: 'Game Art+Animation+Character', labelKey: 'tags.details.characterAnimator' },
              { value: 'Game Art+Animation+Rigging', labelKey: 'tags.details.characterRigger' },
              { value: 'Game Art+Animation+Motion', labelKey: 'tags.details.motionDesigner' }
            ]
          },
          {
            value: 'Game Art+VFX',
            labelKey: 'tags.details.vfxArt',
            children: [
              { value: 'Game Art+VFX+Realtime', labelKey: 'tags.details.realtimeVfxArtist' },
              { value: 'Game Art+VFX+Particles', labelKey: 'tags.details.particleArtist' },
              { value: 'Game Art+VFX+Shaders', labelKey: 'tags.details.shaderVfxArtist' }
            ]
          },
          {
            value: 'Game Art+UI',
            labelKey: 'tags.details.gameUiArt',
            children: [
              { value: 'Game Art+UI+Designer', labelKey: 'tags.details.gameUiDesigner' },
              { value: 'Game Art+UI+UX', labelKey: 'tags.details.gameUxDesigner' }
            ]
          }
        ]
      },
      {
        value: 'Game Audio',
        labelKey: 'tags.specializations.gameAudio',
        children: [
          { value: 'Game Audio+Sound Design', labelKey: 'tags.details.soundDesign' },
          { value: 'Game Audio+Music', labelKey: 'tags.details.gameMusic' },
          { value: 'Game Audio+Technical', labelKey: 'tags.details.technicalAudio' }
        ]
      },
      {
        value: 'Technical Art',
        labelKey: 'tags.specializations.technicalArt',
        children: [
          { value: 'Technical Art+Shaders', labelKey: 'tags.details.shaderArt' },
          { value: 'Technical Art+Pipeline', labelKey: 'tags.details.pipelineArt' },
          { value: 'Technical Art+Tools', labelKey: 'tags.details.toolsArt' }
        ]
      }
    ]
  },

  // ===== STARTUP =====
  {
    value: 'Startup',
    labelKey: 'tags.categories.startup',
    searchTerms: ['стартап', 'предприниматель'],
    children: [
      {
        value: 'Startup+Development',
        labelKey: 'tags.subcategories.startupDevelopment',
        searchTerms: ['разработка'],
        children: [
          {
            value: 'Startup+Development+Backend',
            labelKey: 'tags.specializations.startupBackend',
            children: [
              { value: 'Startup+Development+Backend+Go', labelKey: 'tags.details.startupGoDeveloper' },
              { value: 'Startup+Development+Backend+Node.js', labelKey: 'tags.details.startupNodeDeveloper' },
              { value: 'Startup+Development+Backend+Python', labelKey: 'tags.details.startupPythonDeveloper' }
            ]
          },
          {
            value: 'Startup+Development+Frontend',
            labelKey: 'tags.specializations.startupFrontend',
            children: [
              { value: 'Startup+Development+Frontend+React', labelKey: 'tags.details.startupReactDeveloper' },
              { value: 'Startup+Development+Frontend+Vue', labelKey: 'tags.details.startupVueDeveloper' }
            ]
          },
          {
            value: 'Startup+Development+Mobile',
            labelKey: 'tags.specializations.startupMobile',
            children: [
              { value: 'Startup+Development+Mobile+iOS', labelKey: 'tags.details.startupIosDeveloper' },
              { value: 'Startup+Development+Mobile+Android', labelKey: 'tags.details.startupAndroidDeveloper' }
            ]
          }
        ]
      },
      {
        value: 'Startup+ClientWork',
        labelKey: 'tags.subcategories.startupClientWork',
        searchTerms: ['работа с клиентами', 'customer'],
        children: [
          {
            value: 'Startup+ClientWork+Support',
            labelKey: 'tags.specializations.startupSupport',
            children: [
              { value: 'Startup+ClientWork+Support+ChatManager', labelKey: 'tags.details.startupChatManager' },
              { value: 'Startup+ClientWork+Support+EmailManager', labelKey: 'tags.details.startupEmailManager' },
              { value: 'Startup+ClientWork+Support+VoiceSupport', labelKey: 'tags.details.startupVoiceSupport' }
            ]
          },
          {
            value: 'Startup+ClientWork+Success',
            labelKey: 'tags.specializations.startupSuccess',
            children: [
              { value: 'Startup+ClientWork+Success+CSM', labelKey: 'tags.details.startupCsm' },
              { value: 'Startup+ClientWork+Success+AccountManager', labelKey: 'tags.details.startupAccountManager' }
            ]
          }
        ]
      },
      {
        value: 'Startup+Product',
        labelKey: 'tags.subcategories.startupProduct',
        children: [
          {
            value: 'Startup+Product+Management',
            labelKey: 'tags.specializations.startupProductManagement',
            children: [
              { value: 'Startup+Product+Management+PM', labelKey: 'tags.details.startupProductManager' },
              { value: 'Startup+Product+Management+PO', labelKey: 'tags.details.startupProductOwner' },
              { value: 'Startup+Product+Management+Analyst', labelKey: 'tags.details.startupProductAnalyst' }
            ]
          }
        ]
      },
      {
        value: 'Startup+Growth',
        labelKey: 'tags.subcategories.startupGrowth',
        children: [
          {
            value: 'Startup+Growth+Marketing',
            labelKey: 'tags.specializations.startupGrowthMarketing',
            children: [
              { value: 'Startup+Growth+Marketing+Manager', labelKey: 'tags.details.startupGrowthMarketingManager' },
              { value: 'Startup+Growth+Marketing+Specialist', labelKey: 'tags.details.startupGrowthMarketingSpecialist' }
            ]
          },
          {
            value: 'Startup+Growth+Sales',
            labelKey: 'tags.specializations.startupGrowthSales',
            children: [
              { value: 'Startup+Growth+Sales+Manager', labelKey: 'tags.details.startupSalesManager' },
              { value: 'Startup+Growth+Sales+Representative', labelKey: 'tags.details.startupSalesRep' }
            ]
          }
        ]
      }
    ]
  },
  // ===== SUPPORT =====
  {
    value: 'Support',
    labelKey: 'tags.categories.support',
    searchTerms: ['поддержка', 'служба'],
    children: [
      {
        value: 'Customer Support',
        labelKey: 'tags.specializations.customerSupport',
        children: [
          { value: 'Support+Chat', labelKey: 'tags.details.supportChat' },
          { value: 'Support+Phone', labelKey: 'tags.details.supportPhone' },
          { value: 'Support+Email', labelKey: 'tags.details.supportEmail' },
          { value: 'Support+Level 1', labelKey: 'tags.details.supportL1' },
          { value: 'Support+Level 2', labelKey: 'tags.details.supportL2' },
          { value: 'Support+Level 3', labelKey: 'tags.details.supportL3' }
        ]
      },
      {
        value: 'Technical Support',
        labelKey: 'tags.specializations.technicalSupport',
        children: [
          { value: 'Tech Support+IT', labelKey: 'tags.details.techSupportIt' },
          { value: 'Tech Support+Hardware', labelKey: 'tags.details.techSupportHardware' },
          { value: 'Tech Support+Software', labelKey: 'tags.details.techSupportSoftware' }
        ]
      }
    ]
  },

  // ===== MARKETING & SALES =====
  {
    value: 'Marketing',
    labelKey: 'tags.categories.marketing',
    searchTerms: ['маркетинг', 'реклама'],
    children: [
      {
        value: 'Digital Marketing',
        labelKey: 'tags.specializations.digitalMarketing',
        children: [
          { value: 'Marketing+SEO', labelKey: 'tags.details.marketingSeo' },
          { value: 'Marketing+PPC', labelKey: 'tags.details.marketingPpc' },
          { value: 'Marketing+Content', labelKey: 'tags.details.marketingContent' },
          { value: 'Marketing+SMM', labelKey: 'tags.details.marketingSmm' }
        ]
      },
      {
        value: 'Sales',
        labelKey: 'tags.specializations.sales',
        children: [
          { value: 'Sales+B2B', labelKey: 'tags.details.salesB2b' },
          { value: 'Sales+B2C', labelKey: 'tags.details.salesB2c' },
          { value: 'Sales+Account', labelKey: 'tags.details.salesAccount' }
        ]
      }
    ]
  },

  // ===== MANAGEMENT =====
  {
    value: 'Management',
    labelKey: 'tags.categories.management',
    searchTerms: ['менеджмент', 'управление'],
    children: [
      {
        value: 'Executive',
        labelKey: 'tags.specializations.executive',
        children: [
          { value: 'Management+CEO', labelKey: 'tags.details.ceo' },
          { value: 'Management+CTO', labelKey: 'tags.details.cto' },
          { value: 'Management+CFO', labelKey: 'tags.details.cfo' },
          { value: 'Management+CMO', labelKey: 'tags.details.cmo' },
          { value: 'Management+COO', labelKey: 'tags.details.coo' }
        ]
      },
      {
        value: 'Team Lead',
        labelKey: 'tags.specializations.teamLead',
        children: [
          { value: 'Lead+Engineering', labelKey: 'tags.details.leadEngineering' },
          { value: 'Lead+Product', labelKey: 'tags.details.leadProduct' },
          { value: 'Lead+Design', labelKey: 'tags.details.leadDesign' }
        ]
      },
      {
        value: 'Project Management',
        labelKey: 'tags.specializations.projectManagement',
        children: [
          { value: 'PM+Agile', labelKey: 'tags.details.pmAgile' },
          { value: 'PM+Scrum', labelKey: 'tags.details.pmScrum' },
          { value: 'PM+Waterfall', labelKey: 'tags.details.pmWaterfall' }
        ]
      }
    ]
  },

  // ===== HR & FINANCE =====
  {
    value: 'HR',
    labelKey: 'tags.categories.hr',
    searchTerms: ['HR', 'кадры', 'персонал'],
    children: [
      {
        value: 'Recruitment',
        labelKey: 'tags.specializations.recruitment',
        children: [
          { value: 'HR+Recruiter', labelKey: 'tags.details.recruiter' },
          { value: 'HR+Sourcer', labelKey: 'tags.details.sourcer' },
          { value: 'HR+Tech Recruiter', labelKey: 'tags.details.techRecruiter' }
        ]
      },
      {
        value: 'HR Management',
        labelKey: 'tags.specializations.hrManagement',
        children: [
          { value: 'HR+Manager', labelKey: 'tags.details.hrManager' },
          { value: 'HR+Business Partner', labelKey: 'tags.details.hrBp' }
        ]
      }
    ]
  },

  {
    value: 'Finance',
    labelKey: 'tags.categories.finance',
    searchTerms: ['финансы', 'бухгалтерия'],
    children: [
      {
        value: 'Accounting',
        labelKey: 'tags.specializations.accounting',
        children: [
          { value: 'Finance+Accountant', labelKey: 'tags.details.accountant' },
          { value: 'Finance+Auditor', labelKey: 'tags.details.auditor' }
        ]
      },
      {
        value: 'Financial Analysis',
        labelKey: 'tags.specializations.financialAnalysis',
        children: [
          { value: 'Finance+Analyst', labelKey: 'tags.details.financialAnalyst' },
          { value: 'Finance+Controller', labelKey: 'tags.details.controller' }
        ]
      }
    ]
  },

  // ===== SOFTWARE DEVELOPMENT =====
  {
    value: 'SoftwareDevelopment',
    labelKey: 'tags.categories.softwareDevelopment',
    searchTerms: ['разработка', 'software', 'development', 'программирование'],
    children: [
      {
        value: 'BackendDeveloper',
        labelKey: 'tags.specializations.backendDeveloper',
        children: [
          { value: 'SoftwareDevelopment+Backend+Node.js', labelKey: 'tags.details.sdBackendNodejs' },
          { value: 'SoftwareDevelopment+Backend+Python', labelKey: 'tags.details.sdBackendPython' },
          { value: 'SoftwareDevelopment+Backend+Go', labelKey: 'tags.details.sdBackendGo' },
          { value: 'SoftwareDevelopment+Backend+C#', labelKey: 'tags.details.sdBackendCsharp' },
          { value: 'SoftwareDevelopment+Backend+Java', labelKey: 'tags.details.sdBackendJava' },
          { value: 'SoftwareDevelopment+Backend+PHP', labelKey: 'tags.details.sdBackendPhp' },
          { value: 'SoftwareDevelopment+Backend+Ruby', labelKey: 'tags.details.sdBackendRuby' }
        ]
      },
      {
        value: 'FrontendDeveloper',
        labelKey: 'tags.specializations.frontendDeveloper',
        children: [
          { value: 'SoftwareDevelopment+Frontend+React', labelKey: 'tags.details.sdFrontendReact' },
          { value: 'SoftwareDevelopment+Frontend+Vue', labelKey: 'tags.details.sdFrontendVue' },
          { value: 'SoftwareDevelopment+Frontend+Angular', labelKey: 'tags.details.sdFrontendAngular' }
        ]
      },
      {
        value: 'FullstackDeveloper',
        labelKey: 'tags.specializations.fullstackDeveloper'
      },
      {
        value: 'WebDeveloper',
        labelKey: 'tags.specializations.webDeveloper'
      },
      {
        value: 'ApplicationDeveloper',
        labelKey: 'tags.specializations.applicationDeveloper'
      },
      {
        value: 'MobileApplicationDeveloper',
        labelKey: 'tags.specializations.mobileApplicationDeveloper',
        children: [
          { value: 'SoftwareDevelopment+Mobile+iOS', labelKey: 'tags.details.sdMobileIos' },
          { value: 'SoftwareDevelopment+Mobile+Android', labelKey: 'tags.details.sdMobileAndroid' },
          { value: 'SoftwareDevelopment+Mobile+React Native', labelKey: 'tags.details.sdMobileReactNative' },
          { value: 'SoftwareDevelopment+Mobile+Flutter', labelKey: 'tags.details.sdMobileFlutter' }
        ]
      },
      {
        value: 'ReleaseManager',
        labelKey: 'tags.specializations.releaseManager'
      },
      {
        value: 'GameDeveloper',
        labelKey: 'tags.specializations.gameDeveloper'
      },
      {
        value: 'SoftwareDeveloper',
        labelKey: 'tags.specializations.softwareDeveloper'
      },
      {
        value: 'DatabaseDeveloper',
        labelKey: 'tags.specializations.databaseDeveloper'
      },
      {
        value: 'EmbeddedSoftwareEngineer',
        labelKey: 'tags.specializations.embeddedSoftwareEngineer'
      },
      {
        value: 'HTMLCoding',
        labelKey: 'tags.specializations.htmlCoding'
      },
      {
        value: '1CDeveloper',
        labelKey: 'tags.specializations.1cDeveloper'
      },
      {
        value: 'SoftwareArchitect',
        labelKey: 'tags.specializations.softwareArchitect'
      },
      {
        value: 'SystemSoftwareEngineer',
        labelKey: 'tags.specializations.systemSoftwareEngineer'
      },
      {
        value: 'ERPDeveloper',
        labelKey: 'tags.specializations.erpDeveloper'
      },
      {
        value: 'DatabaseArchitect',
        labelKey: 'tags.specializations.databaseArchitect'
      },
      {
        value: 'HardwareEngineer',
        labelKey: 'tags.specializations.hardwareEngineer'
      },
      {
        value: '1CArchitect',
        labelKey: 'tags.specializations.1cArchitect'
      },
      {
        value: 'OtherDevelopment',
        labelKey: 'tags.specializations.otherDevelopment'
      }
    ]
  },

  // ===== QUALITY ASSURANCE (TESTING) =====
  {
    value: 'QualityAssurance',
    labelKey: 'tags.categories.qualityAssurance',
    searchTerms: ['тестирование', 'qa', 'quality', 'testing'],
    children: [
      {
        value: 'ManualTestEngineer',
        labelKey: 'tags.specializations.manualTestEngineer'
      },
      {
        value: 'TestAutomationEngineer',
        labelKey: 'tags.specializations.testAutomationEngineer'
      },
      {
        value: 'QAEngineer',
        labelKey: 'tags.specializations.qaEngineer'
      },
      {
        value: 'UXTester',
        labelKey: 'tags.specializations.uxTester'
      },
      {
        value: 'SoftwarePerformanceEngineer',
        labelKey: 'tags.specializations.softwarePerformanceEngineer'
      },
      {
        value: 'QAAnalyst',
        labelKey: 'tags.specializations.qaAnalyst'
      },
      {
        value: 'QAManager',
        labelKey: 'tags.specializations.qaManager'
      },
      {
        value: 'QADirector',
        labelKey: 'tags.specializations.qaDirector'
      }
    ]
  },

  // ===== ANALYTICS =====
  {
    value: 'Analytics',
    labelKey: 'tags.categories.analytics',
    searchTerms: ['аналитика', 'analytics', 'analyst', 'данные', 'data'],
    children: [
      {
        value: 'SystemsAnalyst',
        labelKey: 'tags.specializations.systemsAnalyst'
      },
      {
        value: 'BusinessAnalyst',
        labelKey: 'tags.specializations.businessAnalyst'
      },
      {
        value: 'DataAnalyst',
        labelKey: 'tags.specializations.dataAnalyst'
      },
      {
        value: 'UXAnalyst',
        labelKey: 'tags.specializations.uxAnalyst'
      },
      {
        value: 'GameAnalyst',
        labelKey: 'tags.specializations.gameAnalyst'
      },
      {
        value: 'DataEngineer',
        labelKey: 'tags.specializations.dataEngineer'
      },
      {
        value: 'SoftwareAnalyst',
        labelKey: 'tags.specializations.softwareAnalyst'
      },
      {
        value: 'ProductAnalyst',
        labelKey: 'tags.specializations.productAnalyst'
      },
      {
        value: 'BIDeveloper',
        labelKey: 'tags.specializations.biDeveloper'
      },
      {
        value: 'WebAnalyst',
        labelKey: 'tags.specializations.webAnalyst'
      },
      {
        value: '1CAnalyst',
        labelKey: 'tags.specializations.1cAnalyst'
      }
    ]
  },

  // ===== DESIGN =====
  {
    value: 'DesignCategory',
    labelKey: 'tags.categories.designCategory',
    searchTerms: ['дизайн', 'design', 'designer', 'artist'],
    children: [
      {
        value: 'ProductDesigner',
        labelKey: 'tags.specializations.productDesignerSpec'
      },
      {
        value: 'UIUXDesigner',
        labelKey: 'tags.specializations.uiUxDesigner'
      },
      {
        value: 'WebDesigner',
        labelKey: 'tags.specializations.webDesigner'
      },
      {
        value: 'GraphicDesigner',
        labelKey: 'tags.specializations.graphicDesigner'
      },
      {
        value: 'ApplicationDesigner',
        labelKey: 'tags.specializations.applicationDesigner'
      },
      {
        value: 'DesignerIllustrator',
        labelKey: 'tags.specializations.designerIllustrator'
      },
      {
        value: 'GameDesigner',
        labelKey: 'tags.specializations.gameDesignerSpec'
      },
      {
        value: 'NarrativeDesigner',
        labelKey: 'tags.specializations.narrativeDesignerSpec'
      },
      {
        value: 'MotionDesigner',
        labelKey: 'tags.specializations.motionDesigner'
      },
      {
        value: '3DAnimator',
        labelKey: 'tags.specializations.3dAnimator'
      },
      {
        value: 'FlashAnimator',
        labelKey: 'tags.specializations.flashAnimator'
      },
      {
        value: '3DModeler',
        labelKey: 'tags.specializations.3dModeler'
      },
      {
        value: 'ComputerGraphicsArtist',
        labelKey: 'tags.specializations.computerGraphicsArtist'
      },
      {
        value: 'VUIDesigner',
        labelKey: 'tags.specializations.vuiDesigner'
      },
      {
        value: 'ArtDirectorSpec',
        labelKey: 'tags.specializations.artDirectorSpec'
      }
    ]
  },

  // ===== MANAGEMENT (PROJECT & PRODUCT) =====
  {
    value: 'ProjectManagement',
    labelKey: 'tags.categories.projectManagement',
    searchTerms: ['менеджмент', 'management', 'manager', 'проект', 'project'],
    children: [
      {
        value: 'ProjectManager',
        labelKey: 'tags.specializations.projectManagerSpec'
      },
      {
        value: 'ProjectDirector',
        labelKey: 'tags.specializations.projectDirector'
      },
      {
        value: 'ProductManagerSpec',
        labelKey: 'tags.specializations.productManagerSpec'
      },
      {
        value: 'ScrumMaster',
        labelKey: 'tags.specializations.scrumMaster'
      },
      {
        value: 'DeliveryManager',
        labelKey: 'tags.specializations.deliveryManager'
      },
      {
        value: 'CommunityManager',
        labelKey: 'tags.specializations.communityManager'
      },
      {
        value: 'ProductMarketingManager',
        labelKey: 'tags.specializations.productMarketingManager'
      },
      {
        value: 'ProgramManager',
        labelKey: 'tags.specializations.programManager'
      }
    ]
  },

  // ===== INFORMATION SECURITY =====
  {
    value: 'InformationSecurity',
    labelKey: 'tags.categories.informationSecurity',
    searchTerms: ['безопасность', 'security', 'infosec', 'information security'],
    children: [
      {
        value: 'Pentester',
        labelKey: 'tags.specializations.pentester'
      },
      {
        value: 'SecurityAdministrator',
        labelKey: 'tags.specializations.securityAdministrator'
      },
      {
        value: 'SOCAnalyst',
        labelKey: 'tags.specializations.socAnalyst'
      },
      {
        value: 'InformationSecuritySpecialist',
        labelKey: 'tags.specializations.informationSecuritySpecialist'
      },
      {
        value: 'ReverseEngineer',
        labelKey: 'tags.specializations.reverseEngineer'
      },
      {
        value: 'AppSecEngineer',
        labelKey: 'tags.specializations.appSecEngineer'
      },
      {
        value: 'SecurityEngineer',
        labelKey: 'tags.specializations.securityEngineer'
      },
      {
        value: 'NLPEngineer',
        labelKey: 'tags.specializations.nlpEngineer'
      },
      {
        value: 'AntifraudAnalyst',
        labelKey: 'tags.specializations.antifraudAnalyst'
      },
      {
        value: 'InformationSecurityArchitect',
        labelKey: 'tags.specializations.informationSecurityArchitect'
      }
    ]
  },

  // ===== TOP MANAGEMENT =====
  {
    value: 'TopManagement',
    labelKey: 'tags.categories.topManagement',
    searchTerms: ['топ-менеджмент', 'top management', 'ceo', 'cto', 'cfo', 'директор'],
    children: [
      {
        value: 'CEO',
        labelKey: 'tags.specializations.ceoSpec'
      },
      {
        value: 'CTO',
        labelKey: 'tags.specializations.ctoSpec'
      },
      {
        value: 'CFO',
        labelKey: 'tags.specializations.cfoSpec'
      },
      {
        value: 'CMO',
        labelKey: 'tags.specializations.cmoSpec'
      },
      {
        value: 'COO',
        labelKey: 'tags.specializations.cooSpec'
      }
    ]
  },

  // ===== ARTIFICIAL INTELLIGENCE =====
  {
    value: 'ArtificialIntelligence',
    labelKey: 'tags.categories.artificialIntelligence',
    searchTerms: ['искусственный интеллект', 'ai', 'artificial intelligence', 'ml', 'machine learning'],
    children: [
      {
        value: 'DataScientist',
        labelKey: 'tags.specializations.dataScientistSpec'
      },
      {
        value: 'MLEngineer',
        labelKey: 'tags.specializations.mlEngineer'
      }
    ]
  },
  // ===== SUPPORT =====
  {
    value: 'SupportCategory',
    labelKey: 'tags.categories.supportCategory',
    searchTerms: ['поддержка', 'support', 'customer support', 'техподдержка'],
    children: [
      {
        value: 'CustomerSupportSpec',
        labelKey: 'tags.specializations.customerSupportSpec'
      },
      {
        value: 'TechnicalSupportSpec',
        labelKey: 'tags.specializations.technicalSupportSpec'
      },
      {
        value: 'HelpDeskSupport',
        labelKey: 'tags.specializations.helpDeskSupport'
      }
    ]
  },

  // ===== MANUFACTURING AND CONSTRUCTION =====
  {
    value: 'ManufacturingAndConstruction',
    labelKey: 'tags.categories.manufacturingAndConstruction',
    searchTerms: ['производство', 'строительство', 'manufacturing', 'construction'],
    children: [
      {
        value: 'ProductionEngineer',
        labelKey: 'tags.specializations.productionEngineer'
      },
      {
        value: 'QualityEngineerSpec',
        labelKey: 'tags.specializations.qualityEngineerSpec'
      },
      {
        value: 'SafetyEngineerSpec',
        labelKey: 'tags.specializations.safetyEngineerSpec'
      },
      {
        value: 'ConstructionManager',
        labelKey: 'tags.specializations.constructionManager'
      },
      {
        value: 'Foreman',
        labelKey: 'tags.specializations.foreman'
      }
    ]
  },

  // ===== MARKETING =====
  {
    value: 'MarketingCategory',
    labelKey: 'tags.categories.marketingCategory',
    searchTerms: ['маркетинг', 'marketing', 'продвижение', 'реклама'],
    children: [
      {
        value: 'DigitalMarketingSpec',
        labelKey: 'tags.specializations.digitalMarketingSpec'
      },
      {
        value: 'ContentMarketingSpec',
        labelKey: 'tags.specializations.contentMarketingSpec'
      },
      {
        value: 'SMMManagerSpec',
        labelKey: 'tags.specializations.smmManagerSpec'
      },
      {
        value: 'SEOSpecialist',
        labelKey: 'tags.specializations.seoSpecialist'
      },
      {
        value: 'PPCSpecialist',
        labelKey: 'tags.specializations.ppcSpecialist'
      },
      {
        value: 'EmailMarketingSpec',
        labelKey: 'tags.specializations.emailMarketingSpec'
      },
      {
        value: 'BrandManager',
        labelKey: 'tags.specializations.brandManager'
      }
    ]
  },

  // ===== ADMINISTRATION =====
  {
    value: 'AdministrationCategory',
    labelKey: 'tags.categories.administrationCategory',
    searchTerms: ['администрирование', 'administration', 'системное администрирование'],
    children: [
      {
        value: 'SystemAdministrator',
        labelKey: 'tags.specializations.systemAdministrator'
      },
      {
        value: 'NetworkAdministrator',
        labelKey: 'tags.specializations.networkAdministrator'
      },
      {
        value: 'DatabaseAdministrator',
        labelKey: 'tags.specializations.databaseAdministrator'
      },
      {
        value: 'CloudAdministrator',
        labelKey: 'tags.specializations.cloudAdministrator'
      }
    ]
  },

  // ===== CONTENT =====
  {
    value: 'ContentCategory',
    labelKey: 'tags.categories.contentCategory',
    searchTerms: ['контент', 'content', 'копирайтинг', 'writing'],
    children: [
      {
        value: 'ContentManager',
        labelKey: 'tags.specializations.contentManager'
      },
      {
        value: 'CopywriterSpec',
        labelKey: 'tags.specializations.copywriterSpec'
      },
      {
        value: 'TechnicalWriterSpec',
        labelKey: 'tags.specializations.technicalWriterSpec'
      },
      {
        value: 'ContentStrategist',
        labelKey: 'tags.specializations.contentStrategist'
      },
      {
        value: 'EditorSpec',
        labelKey: 'tags.specializations.editorSpec'
      }
    ]
  },

  // ===== SALES =====
  {
    value: 'SalesCategory',
    labelKey: 'tags.categories.salesCategory',
    searchTerms: ['продажи', 'sales', 'сейлз'],
    children: [
      {
        value: 'SalesManagerSpec',
        labelKey: 'tags.specializations.salesManagerSpec'
      },
      {
        value: 'AccountManagerSpec',
        labelKey: 'tags.specializations.accountManagerSpec'
      },
      {
        value: 'SalesRepresentative',
        labelKey: 'tags.specializations.salesRepresentative'
      },
      {
        value: 'BusinessDevelopmentManager',
        labelKey: 'tags.specializations.businessDevelopmentManager'
      },
      {
        value: 'B2BSalesManager',
        labelKey: 'tags.specializations.b2bSalesManager'
      },
      {
        value: 'B2CSalesManager',
        labelKey: 'tags.specializations.b2cSalesManager'
      }
    ]
  },

  // ===== HR (HUMAN RESOURCES) =====
  {
    value: 'HumanResources',
    labelKey: 'tags.categories.humanResources',
    searchTerms: ['hr', 'human resources', 'кадры', 'персонал', 'рекрутинг'],
    children: [
      {
        value: 'HRManagerSpec',
        labelKey: 'tags.specializations.hrManagerSpec'
      },
      {
        value: 'RecruiterSpec',
        labelKey: 'tags.specializations.recruiterSpec'
      },
      {
        value: 'TalentAcquisitionSpec',
        labelKey: 'tags.specializations.talentAcquisitionSpec'
      },
      {
        value: 'HRBusinessPartner',
        labelKey: 'tags.specializations.hrBusinessPartner'
      },
      {
        value: 'CompensationBenefitsSpec',
        labelKey: 'tags.specializations.compensationBenefitsSpec'
      },
      {
        value: 'LearningDevelopmentSpec',
        labelKey: 'tags.specializations.learningDevelopmentSpec'
      }
    ]
  },

  // ===== OFFICE =====
  {
    value: 'OfficeCategory',
    labelKey: 'tags.categories.officeCategory',
    searchTerms: ['офис', 'office', 'administrative', 'administration'],
    children: [
      {
        value: 'OfficeManagerSpec',
        labelKey: 'tags.specializations.officeManagerSpec'
      },
      {
        value: 'ExecutiveAssistant',
        labelKey: 'tags.specializations.executiveAssistant'
      },
      {
        value: 'Receptionist',
        labelKey: 'tags.specializations.receptionist'
      },
      {
        value: 'AdministrativeAssistant',
        labelKey: 'tags.specializations.administrativeAssistant'
      },
      {
        value: 'OfficeCoordinator',
        labelKey: 'tags.specializations.officeCoordinator'
      }
    ]
  },

  // ===== OTHER (ДРУГОЕ) =====
  {
    value: 'Other',
    labelKey: 'tags.categories.other',
    searchTerms: ['другое', 'other', 'прочее', 'остальное'],
    children: [
      { value: 'Remote', labelKey: 'tags.options.remote' },
      { value: 'On-site', labelKey: 'tags.options.onSite' },
      {
        value: 'OtherSupport',
        labelKey: 'tags.specializations.otherSupport'
      },
      {
        value: 'OtherMarketing',
        labelKey: 'tags.specializations.otherMarketing'
      },
      {
        value: 'OtherManagement',
        labelKey: 'tags.specializations.otherManagement'
      },
      {
        value: 'OtherHR',
        labelKey: 'tags.specializations.otherHR'
      },
      {
        value: 'OtherFinance',
        labelKey: 'tags.specializations.otherFinance'
      },
      {
        value: 'OtherBusiness',
        labelKey: 'tags.specializations.otherBusiness'
      },
      {
        value: 'OtherSecurity',
        labelKey: 'tags.specializations.otherSecurity'
      },
      {
        value: 'OtherManufacturing',
        labelKey: 'tags.specializations.otherManufacturing'
      },
      {
        value: 'OtherOffice',
        labelKey: 'tags.specializations.otherOffice'
      }
    ]
  }
]

export const tagsGroups: TagGroup[] = [
  // ===== ОСНОВНЫЕ КАТЕГОРИИ (приоритет 1) =====
  {
    key: 'primary',
    labelKey: 'tags.groups.primary',
    priority: 1,
    category: 'primary',
    options: [
      { value: 'IT', labelKey: 'tags.categories.it', category: 'IT' },
      { value: 'GameDev', labelKey: 'tags.categories.gamedev', category: 'GameDev' },
      { value: 'Startup', labelKey: 'tags.categories.startup', category: 'Startup' },
      { value: 'Other', labelKey: 'tags.categories.other', category: 'Other' },
      { value: 'Support', labelKey: 'tags.categories.support', category: 'Support' },
      { value: 'Marketing', labelKey: 'tags.categories.marketing', category: 'Marketing' },
      { value: 'Management', labelKey: 'tags.categories.management', category: 'Management' },
      { value: 'HR', labelKey: 'tags.categories.hr', category: 'HR' },
      { value: 'Finance', labelKey: 'tags.categories.finance', category: 'Finance' },
      { value: 'Business', labelKey: 'tags.categories.business', category: 'Business' },
      { value: 'InformationSecurity', labelKey: 'tags.categories.informationSecurity', category: 'InformationSecurity' },
      { value: 'ArtificialIntelligence', labelKey: 'tags.categories.artificialIntelligence', category: 'ArtificialIntelligence' },
      { value: 'ManufacturingAndConstruction', labelKey: 'tags.categories.manufacturingAndConstruction', category: 'ManufacturingAndConstruction' }
    ]
  },
  // ===== IT СПЕЦИАЛИЗАЦИИ =====
  {
    key: 'itSpecializations',
    labelKey: 'tags.groups.itSpecializations',
    priority: 2,
    category: 'secondary',
    options: [
      { value: 'Backend', labelKey: 'tags.specializations.backend', parent: 'IT', category: 'IT' },
      { value: 'Frontend', labelKey: 'tags.specializations.frontend', parent: 'IT', category: 'IT' },
      { value: 'Fullstack', labelKey: 'tags.specializations.fullstack', parent: 'IT', category: 'IT' },
      { value: 'Mobile', labelKey: 'tags.specializations.mobile', parent: 'IT', category: 'IT' },
      { value: 'DevOps', labelKey: 'tags.specializations.devops', parent: 'IT', category: 'IT' },
      { value: 'Data', labelKey: 'tags.specializations.data', parent: 'IT', category: 'IT' },
      { value: 'QA', labelKey: 'tags.specializations.qa', parent: 'IT', category: 'IT' },
      { value: 'Security', labelKey: 'tags.specializations.security', parent: 'IT', category: 'IT' },
      { value: 'Design', labelKey: 'tags.specializations.design', parent: 'IT', category: 'IT' }
    ]
  },

  // ===== IT ДЕТАЛЬНЫЕ ТЕГИ =====
  {
    key: 'itDetails',
    labelKey: 'tags.groups.itDetails',
    priority: 3,
    category: 'secondary',
    options: [
      // Backend детали
      { value: 'Backend+Node.js', labelKey: 'tags.details.backendNodejs', parent: 'Backend', category: 'IT' },
      { value: 'Backend+Python', labelKey: 'tags.details.backendPython', parent: 'Backend', category: 'IT' },
      { value: 'Backend+Go', labelKey: 'tags.details.backendGo', parent: 'Backend', category: 'IT' },
      { value: 'Backend+C#', labelKey: 'tags.details.backendCsharp', parent: 'Backend', category: 'IT' },
      { value: 'Backend+Java', labelKey: 'tags.details.backendJava', parent: 'Backend', category: 'IT' },
      { value: 'Backend+PHP', labelKey: 'tags.details.backendPhp', parent: 'Backend', category: 'IT' },
      { value: 'Backend+Ruby', labelKey: 'tags.details.backendRuby', parent: 'Backend', category: 'IT' },
      // Frontend детали
      { value: 'Frontend+React', labelKey: 'tags.details.frontendReact', parent: 'Frontend', category: 'IT' },
      { value: 'Frontend+Vue', labelKey: 'tags.details.frontendVue', parent: 'Frontend', category: 'IT' },
      { value: 'Frontend+Angular', labelKey: 'tags.details.frontendAngular', parent: 'Frontend', category: 'IT' },
      { value: 'Frontend+Next.js', labelKey: 'tags.details.frontendNextjs', parent: 'Frontend', category: 'IT' },
      { value: 'Frontend+Nuxt', labelKey: 'tags.details.frontendNuxt', parent: 'Frontend', category: 'IT' },
      // Mobile детали
      { value: 'Mobile+iOS', labelKey: 'tags.details.mobileIos', parent: 'Mobile', category: 'IT' },
      { value: 'Mobile+Android', labelKey: 'tags.details.mobileAndroid', parent: 'Mobile', category: 'IT' },
      { value: 'Mobile+React Native', labelKey: 'tags.details.mobileReactNative', parent: 'Mobile', category: 'IT' },
      { value: 'Mobile+Flutter', labelKey: 'tags.details.mobileFlutter', parent: 'Mobile', category: 'IT' },
      // DevOps детали
      { value: 'DevOps+AWS', labelKey: 'tags.details.devopsAws', parent: 'DevOps', category: 'IT' },
      { value: 'DevOps+Azure', labelKey: 'tags.details.devopsAzure', parent: 'DevOps', category: 'IT' },
      { value: 'DevOps+GCP', labelKey: 'tags.details.devopsGcp', parent: 'DevOps', category: 'IT' },
      { value: 'DevOps+Kubernetes', labelKey: 'tags.details.devopsK8s', parent: 'DevOps', category: 'IT' },
      { value: 'DevOps+Docker', labelKey: 'tags.details.devopsDocker', parent: 'DevOps', category: 'IT' },
      // Data детали
      { value: 'Data+Analytics', labelKey: 'tags.details.dataAnalytics', parent: 'Data', category: 'IT' },
      { value: 'Data+Science', labelKey: 'tags.details.dataScience', parent: 'Data', category: 'IT' },
      { value: 'Data+Engineering', labelKey: 'tags.details.dataEngineering', parent: 'Data', category: 'IT' },
      { value: 'Data+ML/AI', labelKey: 'tags.details.dataMlAi', parent: 'Data', category: 'IT' },
      // QA детали
      { value: 'QA+Manual', labelKey: 'tags.details.qaManual', parent: 'QA', category: 'IT' },
      { value: 'QA+Automation', labelKey: 'tags.details.qaAutomation', parent: 'QA', category: 'IT' },
      { value: 'QA+Performance', labelKey: 'tags.details.qaPerformance', parent: 'QA', category: 'IT' },
      // Design детали
      { value: 'Design+UX', labelKey: 'tags.details.designUx', parent: 'Design', category: 'IT' },
      { value: 'Design+UI', labelKey: 'tags.details.designUi', parent: 'Design', category: 'IT' },
      { value: 'Design+Product', labelKey: 'tags.details.designProduct', parent: 'Design', category: 'IT' }
    ]
  },

  // ===== GAMEDEV СПЕЦИАЛИЗАЦИИ =====
  {
    key: 'gamedevSpecializations',
    labelKey: 'tags.groups.gamedevSpecializations',
    priority: 2,
    category: 'secondary',
    options: [
      { value: 'Unity', labelKey: 'tags.specializations.unity', parent: 'GameDev', category: 'GameDev' },
      { value: 'Unreal Engine', labelKey: 'tags.specializations.unreal', parent: 'GameDev', category: 'GameDev' },
      { value: 'Game Design', labelKey: 'tags.specializations.gameDesign', parent: 'GameDev', category: 'GameDev' },
      { value: 'Game Art', labelKey: 'tags.specializations.gameArt', parent: 'GameDev', category: 'GameDev' },
      { value: 'Game Audio', labelKey: 'tags.specializations.gameAudio', parent: 'GameDev', category: 'GameDev' },
      { value: 'Technical Art', labelKey: 'tags.specializations.technicalArt', parent: 'GameDev', category: 'GameDev' }
    ]
  },

  // ===== GAMEDEV ДЕТАЛЬНЫЕ ТЕГИ =====
  {
    key: 'gamedevDetails',
    labelKey: 'tags.groups.gamedevDetails',
    priority: 3,
    category: 'secondary',
    options: [
      // Unity детали
      { value: 'Unity+C#', labelKey: 'tags.details.unityCsharp', parent: 'Unity', category: 'GameDev' },
      { value: 'Unity+Mobile', labelKey: 'tags.details.unityMobile', parent: 'Unity', category: 'GameDev' },
      { value: 'Unity+PC', labelKey: 'tags.details.unityPc', parent: 'Unity', category: 'GameDev' },
      // Unreal Engine детали
      { value: 'UE5+C++', labelKey: 'tags.details.ue5Cpp', parent: 'Unreal Engine', category: 'GameDev' },
      { value: 'UE5+C#', labelKey: 'tags.details.ue5Csharp', parent: 'Unreal Engine', category: 'GameDev' },
      { value: 'UE5+Blueprints', labelKey: 'tags.details.ue5Blueprints', parent: 'Unreal Engine', category: 'GameDev' },
      // Game Design детали
      { value: 'Game Design+Level', labelKey: 'tags.details.levelDesign', parent: 'Game Design', category: 'GameDev' },
      { value: 'Game Design+Narrative', labelKey: 'tags.details.narrativeDesign', parent: 'Game Design', category: 'GameDev' },
      { value: 'Game Design+Systems', labelKey: 'tags.details.systemsDesign', parent: 'Game Design', category: 'GameDev' },
      // Game Art детали
      { value: 'Game Art+3D', labelKey: 'tags.details.gameArt3d', parent: 'Game Art', category: 'GameDev' },
      { value: 'Game Art+2D', labelKey: 'tags.details.gameArt2d', parent: 'Game Art', category: 'GameDev' },
      { value: 'Game Art+Character', labelKey: 'tags.details.characterArt', parent: 'Game Art', category: 'GameDev' },
      { value: 'Game Art+Environment', labelKey: 'tags.details.environmentArt', parent: 'Game Art', category: 'GameDev' },
      { value: 'Game Art+VFX', labelKey: 'tags.details.vfxArt', parent: 'Game Art', category: 'GameDev' },
      { value: 'Game Art+UI', labelKey: 'tags.details.gameUiArt', parent: 'Game Art', category: 'GameDev' },
      // Game Audio детали
      { value: 'Game Audio+Sound Design', labelKey: 'tags.details.soundDesign', parent: 'Game Audio', category: 'GameDev' },
      { value: 'Game Audio+Music', labelKey: 'tags.details.gameMusic', parent: 'Game Audio', category: 'GameDev' },
      { value: 'Game Audio+Technical', labelKey: 'tags.details.technicalAudio', parent: 'Game Audio', category: 'GameDev' }
    ]
  },

  // ===== STARTUP СПЕЦИАЛИЗАЦИИ =====
  {
    key: 'startupSpecializations',
    labelKey: 'tags.groups.startupSpecializations',
    priority: 2,
    category: 'secondary',
    options: [
      { value: 'Product', labelKey: 'tags.specializations.product', parent: 'Startup', category: 'Startup' },
      { value: 'Growth', labelKey: 'tags.specializations.growth', parent: 'Startup', category: 'Startup' },
      { value: 'Operations', labelKey: 'tags.specializations.operations', parent: 'Startup', category: 'Startup' }
    ]
  },

  // ===== СТАДИЯ СТАРТАПА =====
  {
    key: 'stage',
    labelKey: 'tags.groups.stage',
    options: [
      { value: 'MVP', labelKey: 'tags.options.mvp' },
      { value: 'Pre-seed', labelKey: 'tags.options.preSeed' },
      { value: 'Seed', labelKey: 'tags.options.seed' },
      { value: 'Series A', labelKey: 'tags.options.seriesA' },
      { value: '0-1', labelKey: 'tags.options.zeroToOne' },
      { value: 'Early-stage', labelKey: 'tags.options.earlyStage' },
      { value: 'Bootstrapped', labelKey: 'tags.options.bootstrapped' },
      { value: 'Venture-backed', labelKey: 'tags.options.ventureBacked' },
      { value: 'Remote', labelKey: 'tags.options.remote' },
      { value: 'Hybrid', labelKey: 'tags.options.hybrid' },
      { value: 'On-site', labelKey: 'tags.options.onSite' },
      { value: 'Part-time', labelKey: 'tags.options.partTime' },
      { value: 'Full-time', labelKey: 'tags.options.fullTime' },
      { value: 'Contract', labelKey: 'tags.options.contract' }
    ]
  },
  // Общие роли
  {
    key: 'role',
    labelKey: 'tags.groups.role',
    options: [
      { value: 'Fullstack', labelKey: 'tags.options.fullstack' },
      { value: 'Backend', labelKey: 'tags.options.backend' },
      { value: 'Frontend', labelKey: 'tags.options.frontend' },
      { value: 'Mobile', labelKey: 'tags.options.mobile' },
      { value: 'Product Manager', labelKey: 'tags.options.product' },
      { value: 'Product Designer', labelKey: 'tags.options.productDesigner' },
      { value: 'UX/UI Designer', labelKey: 'tags.options.uxui' },
      { value: 'Growth', labelKey: 'tags.options.growth' },
      { value: 'Marketing', labelKey: 'tags.options.marketing' },
      { value: 'SMM', labelKey: 'tags.options.smm' },
      { value: 'Sales', labelKey: 'tags.options.sales' },
      { value: 'Operations', labelKey: 'tags.options.operations' },
      { value: 'Finance', labelKey: 'tags.options.finance' },
      { value: 'HR', labelKey: 'tags.options.hr' },
      { value: 'Community Manager', labelKey: 'tags.options.communityManager' },
      { value: 'Customer Support', labelKey: 'tags.options.customerSupport' },
      { value: 'Data Analyst', labelKey: 'tags.options.data' },
      { value: 'DevOps', labelKey: 'tags.options.devops' },
      { value: 'QA', labelKey: 'tags.options.qa' },
      { value: 'Sales Manager', labelKey: 'tags.options.salesManager' }
    ]
  },
  // Геймдев роли
  {
    key: 'gamedevRole',
    labelKey: 'tags.groups.gamedevRole',
    options: [
      { value: 'Game Designer', labelKey: 'tags.options.gameDesigner' },
      { value: 'Level Designer', labelKey: 'tags.options.levelDesigner' },
      { value: 'Narrative Designer', labelKey: 'tags.options.narrativeDesigner' },
      { value: 'Game Producer', labelKey: 'tags.options.gameProducer' },
      { value: 'Gameplay Programmer', labelKey: 'tags.options.gameplayProgrammer' },
      { value: 'Engine Programmer', labelKey: 'tags.options.engineProgrammer' },
      { value: 'Graphics Programmer', labelKey: 'tags.options.graphicsProgrammer' },
      { value: '3D Artist', labelKey: 'tags.options.artist3d' },
      { value: '2D Artist', labelKey: 'tags.options.artist2d' },
      { value: 'Character Artist', labelKey: 'tags.options.characterArtist' },
      { value: 'Environment Artist', labelKey: 'tags.options.environmentArtist' },
      { value: 'Concept Artist', labelKey: 'tags.options.conceptArtist' },
      { value: 'UI Artist', labelKey: 'tags.options.uiArtist' },
      { value: 'Animator', labelKey: 'tags.options.animator' },
      { value: 'Technical Artist', labelKey: 'tags.options.technicalArtist' },
      { value: 'VFX Artist', labelKey: 'tags.options.vfxArtist' },
      { value: 'Sound Designer', labelKey: 'tags.options.soundDesigner' },
      { value: 'Composer', labelKey: 'tags.options.composer' },
      { value: 'Audio Engineer', labelKey: 'tags.options.audioEngineer' }
    ]
  },
  // Менеджмент
  {
    key: 'management',
    labelKey: 'tags.groups.management',
    options: [
      { value: 'CEO', labelKey: 'tags.options.ceo' },
      { value: 'CTO', labelKey: 'tags.options.cto' },
      { value: 'CFO', labelKey: 'tags.options.cfo' },
      { value: 'CMO', labelKey: 'tags.options.cmo' },
      { value: 'CIO', labelKey: 'tags.options.cio' },
      { value: 'COO', labelKey: 'tags.options.coo' },
      { value: 'CCO', labelKey: 'tags.options.cco' },
      { value: 'HRD', labelKey: 'tags.options.hrd' },
      { value: 'CLO', labelKey: 'tags.options.clo' },
      { value: 'Production Manager', labelKey: 'tags.options.productionManager' },
      { value: 'Branch Manager', labelKey: 'tags.options.branchManager' },
      { value: 'Project Manager', labelKey: 'tags.options.projectManager' },
      { value: 'Team Lead', labelKey: 'tags.options.teamLead' }
    ]
  },
  // IT роли
  {
    key: 'itRole',
    labelKey: 'tags.groups.itRole',
    options: [
      { value: 'BI Analyst', labelKey: 'tags.options.biAnalyst' },
      { value: 'Business Analyst', labelKey: 'tags.options.businessAnalyst' },
      { value: 'Data Scientist', labelKey: 'tags.options.dataScientist' },
      { value: 'Product Analyst', labelKey: 'tags.options.productAnalyst' },
      { value: 'System Administrator', labelKey: 'tags.options.sysadmin' },
      { value: 'Network Engineer', labelKey: 'tags.options.networkEngineer' },
      { value: 'System Engineer', labelKey: 'tags.options.systemEngineer' },
      { value: 'System Analyst', labelKey: 'tags.options.systemAnalyst' },
      { value: 'InfoSec Specialist', labelKey: 'tags.options.infosecSpecialist' },
      { value: 'Technical Support', labelKey: 'tags.options.techSupport' },
      { value: 'Technical Writer', labelKey: 'tags.options.technicalWriter' }
    ]
  },
  // Административный персонал
  {
    key: 'admin',
    labelKey: 'tags.groups.role',
    options: [
      { value: 'Administrator', labelKey: 'tags.options.administrator' },
      { value: 'Office Manager', labelKey: 'tags.options.officeManager' },
      { value: 'Secretary', labelKey: 'tags.options.secretary' },
      { value: 'PC Operator', labelKey: 'tags.options.pcOperator' },
      { value: 'Office Admin', labelKey: 'tags.options.officeAdmin' },
      { value: 'Translator', labelKey: 'tags.options.translator' },
      { value: 'Courier', labelKey: 'tags.options.courier' },
      { value: 'Facility Manager', labelKey: 'tags.options.facilityManager' }
    ]
  },
  // Автомобильный бизнес
  {
    key: 'automotive',
    labelKey: 'tags.groups.role',
    options: [
      { value: 'Car Washer', labelKey: 'tags.options.carWasher' },
      { value: 'Auto Mechanic', labelKey: 'tags.options.autoMechanic' },
      { value: 'Service Advisor', labelKey: 'tags.options.serviceAdvisor' }
    ]
  },
  // Безопасность
  {
    key: 'security',
    labelKey: 'tags.groups.role',
    options: [
      { value: 'Security Guard', labelKey: 'tags.options.securityGuard' },
      { value: 'Security Specialist', labelKey: 'tags.options.securitySpecialist' }
    ]
  },
  // Маркетинг и продажи
  {
    key: 'sales',
    labelKey: 'tags.groups.sales',
    options: [
      { value: 'Event Manager', labelKey: 'tags.options.eventManager' },
      { value: 'PR Manager', labelKey: 'tags.options.prManager' },
      { value: 'SMM Manager', labelKey: 'tags.options.smmManager' },
      { value: 'Marketing Analyst', labelKey: 'tags.options.marketingAnalyst' },
      { value: 'Internet Marketer', labelKey: 'tags.options.internetMarketer' },
      { value: 'Partner Manager', labelKey: 'tags.options.partnerManager' },
      { value: 'Promoter', labelKey: 'tags.options.promoter' },
      { value: 'Copywriter', labelKey: 'tags.options.copywriter' }
    ]
  },
  // Производство
  {
    key: 'production',
    labelKey: 'tags.groups.production',
    options: [
      { value: 'Quality Engineer', labelKey: 'tags.options.qualityEngineer' },
      { value: 'Safety Engineer', labelKey: 'tags.options.safetyEngineer' },
      { value: 'Design Engineer', labelKey: 'tags.options.designEngineer' },
      { value: 'Power Engineer', labelKey: 'tags.options.powerEngineer' },
      { value: 'Technologist', labelKey: 'tags.options.technologist' },
      { value: 'Machinist', labelKey: 'tags.options.machinist' },
      { value: 'Mechanic', labelKey: 'tags.options.mechanic' },
      { value: 'CNC Operator', labelKey: 'tags.options.cncOperator' },
      { value: 'Welder', labelKey: 'tags.options.welder' },
      { value: 'Electrician', labelKey: 'tags.options.electrician' },
      { value: 'Plumber', labelKey: 'tags.options.plumber' }
    ]
  },
  // Рабочий персонал
  {
    key: 'labor',
    labelKey: 'tags.groups.role',
    options: [
      { value: 'Loader', labelKey: 'tags.options.loader' },
      { value: 'Warehouse Keeper', labelKey: 'tags.options.warehouseKeeper' },
      { value: 'Packer', labelKey: 'tags.options.packer' },
      { value: 'Laborer', labelKey: 'tags.options.laborer' }
    ]
  },
  // Медицина
  {
    key: 'medical',
    labelKey: 'tags.groups.medical',
    options: [
      { value: 'Doctor', labelKey: 'tags.options.doctor' },
      { value: 'Nurse', labelKey: 'tags.options.nurse' },
      { value: 'Pharmacist', labelKey: 'tags.options.pharmacist' },
      { value: 'Veterinarian', labelKey: 'tags.options.veterinarian' },
      { value: 'Medical Rep', labelKey: 'tags.options.medicalRep' }
    ]
  },
  // Образование
  {
    key: 'education',
    labelKey: 'tags.groups.education',
    options: [
      { value: 'Teacher', labelKey: 'tags.options.teacher' },
      { value: 'Business Trainer', labelKey: 'tags.options.businessTrainer' },
      { value: 'Psychologist', labelKey: 'tags.options.psychologist' }
    ]
  },
  // Сервис
  {
    key: 'service',
    labelKey: 'tags.groups.service',
    options: [
      { value: 'Waiter', labelKey: 'tags.options.waiter' },
      { value: 'Chef', labelKey: 'tags.options.chef' },
      { value: 'Driver', labelKey: 'tags.options.driver' },
      { value: 'Cleaner', labelKey: 'tags.options.cleaner' },
      { value: 'Hairdresser', labelKey: 'tags.options.hairdresser' },
      { value: 'Fitness Trainer', labelKey: 'tags.options.fitnessTrainer' }
    ]
  },
  // Розница
  {
    key: 'retail',
    labelKey: 'tags.groups.role',
    options: [
      { value: 'Sales Assistant', labelKey: 'tags.options.salesAssistant' },
      { value: 'Merchandiser', labelKey: 'tags.options.merchandiser' },
      { value: 'Store Manager', labelKey: 'tags.options.storeDirector' }
    ]
  },
  // Финансы
  {
    key: 'finance',
    labelKey: 'tags.groups.finance',
    options: [
      { value: 'Accountant', labelKey: 'tags.options.accountant' },
      { value: 'Auditor', labelKey: 'tags.options.auditor' },
      { value: 'Financial Analyst', labelKey: 'tags.options.financialAnalyst' },
      { value: 'Financial Manager', labelKey: 'tags.options.financialManager' },
      { value: 'Economist', labelKey: 'tags.options.economist' },
      { value: 'Broker', labelKey: 'tags.options.broker' }
    ]
  },
  // Юриспруденция
  {
    key: 'legal',
    labelKey: 'tags.groups.legal',
    options: [
      { value: 'Lawyer', labelKey: 'tags.options.lawyer' },
      { value: 'Legal Counsel', labelKey: 'tags.options.legalCounsel' },
      { value: 'Compliance Manager', labelKey: 'tags.options.complianceManager' }
    ]
  },
  // Навыки и подходы
  {
    key: 'level',
    labelKey: 'tags.groups.level',
    options: [
      { value: 'Generalist', labelKey: 'tags.options.generalist' },
      { value: 'Owner mindset', labelKey: 'tags.options.ownerMindset' },
      { value: 'Hands-on', labelKey: 'tags.options.handsOn' },
      { value: 'Self-starter', labelKey: 'tags.options.selfStarter' },
      { value: 'Fast learner', labelKey: 'tags.options.fastLearner' },
      { value: 'High autonomy', labelKey: 'tags.options.highAutonomy' },
      { value: 'Fast-paced', labelKey: 'tags.options.fastPaced' },
      { value: 'Ambiguity tolerant', labelKey: 'tags.options.ambiguityTolerant' },
      { value: 'Builder', labelKey: 'tags.options.builder' },
      { value: 'Scrappy', labelKey: 'tags.options.scrappy' }
    ]
  },
  // Технологии
  {
    key: 'stack',
    labelKey: 'tags.groups.stack',
    options: [
      { value: 'JavaScript', labelKey: 'tags.options.javascript' },
      { value: 'TypeScript', labelKey: 'tags.options.typescript' },
      { value: 'Node.js', labelKey: 'tags.options.nodejs' },
      { value: 'Go', labelKey: 'tags.options.go' },
      { value: 'Python', labelKey: 'tags.options.python' },
      { value: 'C#', labelKey: 'tags.options.csharp' },
      { value: 'C++', labelKey: 'tags.options.cpp' },
      { value: 'Rust', labelKey: 'tags.options.rust' },
      { value: 'React', labelKey: 'tags.options.react' },
      { value: 'Vue', labelKey: 'tags.options.vue' },
      { value: 'Nuxt', labelKey: 'tags.options.nuxt' },
      { value: 'Next.js', labelKey: 'tags.options.nextjs' },
      { value: 'PostgreSQL', labelKey: 'tags.options.postgresql' },
      { value: 'MongoDB', labelKey: 'tags.options.mongodb' },
      { value: 'Redis', labelKey: 'tags.options.redis' },
      { value: 'Docker', labelKey: 'tags.options.docker' },
      { value: 'Kubernetes', labelKey: 'tags.options.kubernetes' },
      { value: 'CI/CD', labelKey: 'tags.options.cicd' },
      { value: 'AWS', labelKey: 'tags.options.aws' },
      { value: 'Firebase', labelKey: 'tags.options.firebase' }
    ]
  },
  // Геймдев технологии
  {
    key: 'gamedevStack',
    labelKey: 'tags.groups.gamedevStack',
    options: [
      { value: 'Unity', labelKey: 'tags.options.unity' },
      { value: 'Unreal Engine', labelKey: 'tags.options.unreal' },
      { value: 'Godot', labelKey: 'tags.options.godot' },
      { value: 'GameMaker', labelKey: 'tags.options.gamemaker' },
      { value: 'Blender', labelKey: 'tags.options.blender' },
      { value: 'Maya', labelKey: 'tags.options.maya' },
      { value: 'Substance Painter', labelKey: 'tags.options.substancePainter' },
      { value: 'Photoshop', labelKey: 'tags.options.photoshop' },
      { value: 'Figma', labelKey: 'tags.options.figma' },
      { value: 'Spine', labelKey: 'tags.options.spine' },
      { value: 'FMOD', labelKey: 'tags.options.fmod' },
      { value: 'Wwise', labelKey: 'tags.options.wwise' }
    ]
  },
  // Платформы
  {
    key: 'platforms',
    labelKey: 'tags.groups.gamedevStack',
    options: [
      { value: 'PC', labelKey: 'tags.options.pc' },
      { value: 'Console', labelKey: 'tags.options.console' },
      { value: 'Mobile', labelKey: 'tags.options.mobilePlatform' },
      { value: 'Web', labelKey: 'tags.options.web' },
      { value: 'VR/AR', labelKey: 'tags.options.vr' }
    ]
  },
  // Жанры игр
  {
    key: 'gameGenres',
    labelKey: 'tags.groups.gamedevRole',
    options: [
      { value: 'Casual', labelKey: 'tags.options.casual' },
      { value: 'Hypercasual', labelKey: 'tags.options.hypercasual' },
      { value: 'Midcore', labelKey: 'tags.options.midcore' },
      { value: 'Hardcore', labelKey: 'tags.options.hardcore' },
      { value: 'RPG', labelKey: 'tags.options.rpg' },
      { value: 'Strategy', labelKey: 'tags.options.strategy' },
      { value: 'Shooter', labelKey: 'tags.options.shooter' },
      { value: 'Puzzle', labelKey: 'tags.options.puzzle' }
    ]
  },
  // Сферы продукта
  {
    key: 'domain',
    labelKey: 'tags.groups.domain',
    options: [
      { value: 'FinTech', labelKey: 'tags.options.fintech' },
      { value: 'EdTech', labelKey: 'tags.options.edtech' },
      { value: 'HRTech', labelKey: 'tags.options.hrtech' },
      { value: 'GameDev', labelKey: 'tags.options.gamedev' },
      { value: 'SaaS', labelKey: 'tags.options.saas' },
      { value: 'Marketplace', labelKey: 'tags.options.marketplace' },
      { value: 'B2B', labelKey: 'tags.options.b2b' },
      { value: 'B2C', labelKey: 'tags.options.b2c' },
      { value: 'B2B2C', labelKey: 'tags.options.b2b2c' },
      { value: 'AI/ML', labelKey: 'tags.options.aiml' },
      { value: 'HealthTech', labelKey: 'tags.options.healthtech' },
      { value: 'Social Network', labelKey: 'tags.options.socialNetwork' },
      { value: 'E-commerce', labelKey: 'tags.options.ecommerce' },
      { value: 'Crypto/Web3', labelKey: 'tags.options.crypto' },
      { value: 'Automotive', labelKey: 'tags.options.automotive' }
    ]
  },
  // Условия
  {
    key: 'offer',
    labelKey: 'tags.groups.offer',
    options: [
      { value: 'Equity', labelKey: 'tags.options.equity' },
      { value: 'Options', labelKey: 'tags.options.options' },
      { value: 'Bonus', labelKey: 'tags.options.bonus' },
      { value: 'Flexible schedule', labelKey: 'tags.options.flexibleSchedule' },
      { value: 'Async', labelKey: 'tags.options.async' },
      { value: 'Relocation', labelKey: 'tags.options.relocation' },
      { value: 'Visa support', labelKey: 'tags.options.visaSupport' },
      { value: 'Coworking', labelKey: 'tags.options.coworking' },
      { value: 'Education', labelKey: 'tags.options.education' }
    ]
  },
]

export const tagLabelMap = Object.fromEntries(
  tagsGroups.flatMap(group =>
    group.options.map(option => [option.value, option.labelKey])
  )
)

export const allTags = [
  ...new Set(tagsGroups.flatMap(group => group.options.map(option => option.value)))
]

// ===== ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ДЛЯ РАБОТЫ С ИЕРАРХИЕЙ ТЕГОВ =====

/**
 * Получить все дочерние теги для указанного тега
 */
export function getChildTags(tagValue: string): string[] {
  const children: string[] = []
  
  function traverse(nodes: TagHierarchy[]) {
    for (const node of nodes) {
      if (node.value === tagValue && node.children) {
        // Нашли нужный тег, собираем все дочерние
        function collectChildren(childNodes: TagHierarchy[]) {
          for (const child of childNodes) {
            children.push(child.value)
            if (child.children) {
              collectChildren(child.children)
            }
          }
        }
        collectChildren(node.children)
        return
      }
      
      if (node.children) {
        traverse(node.children)
      }
    }
  }
  
  traverse(tagTaxonomy)
  return children
}

/**
 * Получить родительский тег для указанного тега
 */
export function getParentTag(tagValue: string): string | null {
  function findParent(nodes: TagHierarchy[], parent: string | null = null): string | null {
    for (const node of nodes) {
      if (node.value === tagValue) {
        return parent
      }
      if (node.children) {
        const found = findParent(node.children, node.value)
        if (found !== null) return found
      }
    }
    return null
  }
  
  return findParent(tagTaxonomy)
}

/**
 * Получить всю цепочку родительских тегов (от корня до указанного)
 */
export function getTagPath(tagValue: string): string[] {
  const path: string[] = []
  
  function findPath(nodes: TagHierarchy[], currentPath: string[] = []): boolean {
    for (const node of nodes) {
      const newPath = [...currentPath, node.value]
      
      if (node.value === tagValue) {
        path.push(...newPath)
        return true
      }
      
      if (node.children && findPath(node.children, newPath)) {
        return true
      }
    }
    return false
  }
  
  findPath(tagTaxonomy)
  return path
}

/**
 * Проверить, является ли тег категорией верхнего уровня
 */
export function isPrimaryCategory(tagValue: string): boolean {
  return tagTaxonomy.some(node => node.value === tagValue)
}

/**
 * Получить категорию верхнего уровня для тега
 */
export function getPrimaryCategory(tagValue: string): string | null {
  const path = getTagPath(tagValue)
  return path.length > 0 && path[0] ? path[0] : null
}

/**
 * Поиск тегов по ключевым словам с учетом иерархии
 */
export function searchTags(query: string): string[] {
  const lowerQuery = query.toLowerCase()
  const results: string[] = []
  
  function traverse(nodes: TagHierarchy[]) {
    for (const node of nodes) {
      // Проверяем значение тега
      if (node.value.toLowerCase().includes(lowerQuery)) {
        results.push(node.value)
      }
      
      // Проверяем термины для поиска
      if (node.searchTerms) {
        const matchesSearchTerms = node.searchTerms.some(term => 
          term.toLowerCase().includes(lowerQuery)
        )
        if (matchesSearchTerms && !results.includes(node.value)) {
          results.push(node.value)
        }
      }
      
      // Рекурсивно проходим по дочерним элементам
      if (node.children) {
        traverse(node.children)
      }
    }
  }
  
  traverse(tagTaxonomy)
  return results
}

/**
 * Фильтрация вакансий по тегам с учетом иерархии
 * Если выбран родительский тег, показываем все вакансии с дочерними тегами
 */
export function matchesTagFilter(vacancyTags: string[], filterTags: string[]): boolean {
  if (filterTags.length === 0) return true
  
  return filterTags.every(filterTag => {
    // Проверяем точное совпадение
    if (vacancyTags.includes(filterTag)) return true
    
    // Получаем все дочерние теги фильтра
    const childTags = getChildTags(filterTag)
    
    // Проверяем, есть ли хотя бы один дочерний тег в вакансии
    return childTags.some(childTag => vacancyTags.includes(childTag))
  })
}

/**
 * Получить теги по приоритету (сначала основные категории, потом остальные)
 */
export function getTagsByPriority(): TagGroup[] {
  return [...tagsGroups].sort((a, b) => {
    const priorityA = a.priority || 999
    const priorityB = b.priority || 999
    return priorityA - priorityB
  })
}

/**
 * Получить только основные категории (IT, GameDev, Startup)
 */
export function getPrimaryCategories(): TagOption[] {
  const primaryGroup = tagsGroups.find(g => g.key === 'primary')
  return primaryGroup ? primaryGroup.options : []
}

/**
 * Получить специализации для категории
 */
export function getSpecializationsForCategory(category: string): TagOption[] {
  const groupKey = `${category.toLowerCase()}Specializations`
  const group = tagsGroups.find(g => g.key === groupKey)
  return group ? group.options : []
}

/**
 * Получить детальные теги для специализации
 */
export function getDetailsForSpecialization(specialization: string): TagOption[] {
  const category = getPrimaryCategory(specialization)
  if (!category) return []
  
  const groupKey = `${category.toLowerCase()}Details`
  const group = tagsGroups.find(g => g.key === groupKey)
  
  if (!group) return []
  
  return group.options.filter(opt => opt.parent === specialization)
}

/**
 * Форматировать тег для отображения (разбить комбинированные теги)
 */
export function formatTagDisplay(tag: string): string {
  // Если тег содержит '+', показываем только последнюю часть для краткости
  if (tag.includes('+')) {
    const parts = tag.split('+')
    return parts[parts.length - 1] || tag
  }
  return tag
}

/**
 * Получить полное описание тега с иерархией для tooltips
 */
export function getTagFullDescription(tag: string): string {
  const path = getTagPath(tag)
  return path.join(' > ')
}
