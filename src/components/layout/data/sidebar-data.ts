import {
  Boxes,
  HelpCircle,
  LayoutDashboard,
  MapPinned,
  Package,
  Server,
  Settings,
  ShipWheel,
  Users,
} from 'lucide-react'
import { type SidebarData } from '../types'

export const sidebarData: SidebarData = {
  user: {
    name: 'Aurora Operator',
    email: 'admin@aurora.local',
    avatar: '/avatars/shadcn.jpg',
  },
  teams: [
    {
      name: 'Aurora Platform',
      logo: Boxes,
      plan: 'Operations Console',
    },
  ],
  navGroups: [
    {
      title: 'Control',
      items: [
        {
          title: 'Overview',
          url: '/',
          icon: LayoutDashboard,
        },
        {
          title: 'IAM',
          icon: Users,
          items: [
            {
              title: 'User List',
              url: '/users',
            },
            {
              title: 'Role List',
              url: '/roles',
            },
          ],
        },
        {
          title: 'Hypervisor',
          url: '/hypervisor',
          icon: Server,
        },
        {
          title: 'K8s Resource Platform',
          url: '/k8s',
          icon: ShipWheel,
        },
        {
          title: 'Plans',
          url: '/plans',
          icon: Package,
        },
        {
          title: 'Zones',
          url: '/zones',
          icon: MapPinned,
        },
      ],
    },
    {
      title: 'Workspace',
      items: [
        {
          title: 'Settings',
          url: '/settings/account',
          icon: Settings,
        },
        {
          title: 'Help Center',
          url: '/help-center',
          icon: HelpCircle,
        },
      ],
    },
  ],
}
